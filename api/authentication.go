package api

import (
	"crypto/aes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"maid/util"
	"net"
	"strconv"
	"time"

	"github.com/jmhobbs/skip32"
	"golang.org/x/crypto/chacha20"
)

type X19AuthServerConnection struct {
	// connection
	Address       string
	Dial          func(network, addr string) (net.Conn, error)
	connection    net.Conn
	established   bool
	encryptCipher *chacha20.Cipher
	crc32Table    *crc32.Table

	// authentication
	UserToken string
	EntityId  string
}

func (c *X19AuthServerConnection) HasEstablished() bool {
	return c.connection != nil && c.established
}

func (c *X19AuthServerConnection) Disconnect() {
	if c.connection != nil {
		c.connection.Close()
	}
	c.connection = nil
}

// this will block current thread used to receive message
func (c *X19AuthServerConnection) Establish() error {
	conn, err := c.Dial("tcp", c.Address)
	if err != nil {
		return err
	}
	c.connection = conn
	defer func() {
		c.connection.Close()
		c.connection = nil
		c.established = false
	}()

	// do handshake
	remoteKey, err := c.readVariantBytes()
	if err != nil {
		return err
	}
	localKey, err := util.Xor([]byte(c.UserToken), []byte{0xac, 0x24, 0x9c, 0x69, 0xc7, 0x2c, 0xb3, 0xb4, 0x4e, 0xc0, 0xcc, 0x6c, 0x54, 0x3a, 0x81, 0x95})
	if err != nil {
		return err
	}
	{
		handshakeBody := make([]byte, 0x17)
		handshakeBody[0] = 0x15

		// encrypt entity id
		entityId, err := strconv.ParseUint(c.EntityId, 10, 64)
		if err != nil {
			return err
		}
		encryptedId := skip32.Encrypt([10]byte{'S', 'a', 'i', 'n', 't', 'S', 't', 'e', 'v', 'e'}, uint32(entityId))
		i32tob(encryptedId, &handshakeBody, 2)

		block, err := aes.NewCipher(localKey)
		if err != nil {
			return err
		}
		encryptedToken := make([]byte, block.BlockSize())
		block.Encrypt(encryptedToken, remoteKey)

		for i := 0; i < len(encryptedToken); i++ {
			handshakeBody[6+i] = encryptedToken[i]
		}

		_, err = conn.Write(handshakeBody)
		if err != nil {
			return err
		}

		resp, err := c.readVariantBytes()
		if err != nil {
			return err
		}

		if resp[0] != 0x00 {
			return errors.New("failed handshake (" + hex.EncodeToString(resp) + ")")
		}
	}

	c.crc32Table = crc32.MakeTable(crc32.IEEE)
	cipherNonce := []byte{0x31, 0x36, 0x33, 0x20, 0x4e, 0x65, 0x74, 0x45, 0x61, 0x73, 0x65, 0x0a}
	c.encryptCipher, err = chacha20.NewUnauthenticatedCipher(append(localKey, remoteKey...), cipherNonce)
	if err != nil {
		return err
	}
	decryptCipher, err := chacha20.NewUnauthenticatedCipher(append(remoteKey, localKey...), cipherNonce)
	if err != nil {
		return err
	}

	c.established = true

	go func() {
		for c.established {
			// keep alive
			time.Sleep(30 * time.Second)
			c.SendPacket(0, []byte("iamwpf"))
		}
	}()

	for {
		body, err := c.readVariantBytes()
		if err != nil {
			return err
		}
		println(len(body))
		if len(body) < 4 {
			continue
		}
		decryptedBody := make([]byte, len(body))
		decryptCipher.XORKeyStream(decryptedBody, body)

		checksum := make([]byte, 4)
		binary.LittleEndian.PutUint32(checksum, crc32.Checksum(decryptedBody[4:], c.crc32Table))
		for i := 0; i < 4; i++ {
			if checksum[i] != decryptedBody[i] {
				return errors.New("failed to verify checksum")
			}
		}

		packetId := decryptedBody[4]
		decryptedBody = decryptedBody[8:]

		fmt.Printf("%d %s", packetId, hex.EncodeToString(decryptedBody))
	}
}

func (c *X19AuthServerConnection) readVariantBytes() ([]byte, error) {
	var size uint16
	var buf []byte

	// get size
	buf = make([]byte, 2)
	_, err := c.connection.Read(buf)
	if err != nil {
		return nil, err
	}
	size = binary.LittleEndian.Uint16(buf)

	fmt.Printf("size: %d\n", size)

	buf = make([]byte, size)
	_, err = c.connection.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (c *X19AuthServerConnection) SendPacket(id byte, data []byte) error {
	if !c.established {
		return errors.New("connection not established yet")
	}
	body := make([]byte, len(data)+8)
	body[4] = id
	body[5] = 0x88
	body[6] = 0x88
	body[7] = 0x88
	for i := 0; i < len(data); i++ {
		body[i+8] = data[i]
	}

	checksum := make([]byte, 4)
	binary.LittleEndian.PutUint32(checksum, crc32.Checksum(body[4:], c.crc32Table))
	for i := 0; i < 4; i++ {
		body[i] = checksum[i]
	}
	c.encryptCipher.XORKeyStream(body, body)
	c.connection.Write(make([]byte, 2))
	c.connection.Write(body)

	return nil
}

func i32tob(val uint32, arr *[]byte, offset uint32) {
	for i := uint32(0); i < 4; i++ {
		(*arr)[offset+i] = byte((val >> (8 * i)) & 0xff)
	}
}

func GenerateAuthenticationBody(encryptHash, entityId, version, launcherVersion, mods, launchWrapperMD5, gameDataMD5 string) []byte {
	body := make([]byte, 0)

	entityIdInt, _ := strconv.ParseUint(entityId, 10, 64)
	entityIdByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(entityIdByte, entityIdInt)
	body = append(body, entityIdByte[:6]...) // TODO: entity id sum
	body = append(body, 0x00, byte(len(encryptHash)))
	body = append(body, []byte(encryptHash)...)
	body = append(body, byte(len(launcherVersion)))
	body = append(body, []byte(launcherVersion)...)
	body = append(body, byte(len(version)))
	body = append(body, []byte(version)...)
	body = append(body, make([]byte, 20)...) // TODO: md5

	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(len(mods)))
	body = append(body, b...)
	body = append(body, []byte(mods)...)

	body = append(body, 0x02, 0x00, '[', ']')
	body = append(body, byte(len(encryptHash)), 0x00)
	xorHash := make([]byte, len(encryptHash))
	for i := 0; i < len(encryptHash); i++ {
		xorHash[i] = encryptHash[i] ^ 0x5a
	}
	body = append(body, xorHash...)
	body = append(body, 0x07, 'n', 'e', 't', 'e', 'a', 's', 'e')

	return body
}
