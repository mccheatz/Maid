package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"

	mcnet "github.com/Tnze/go-mc/net"
	"github.com/Tnze/go-mc/net/packet"
)

const (
	HANDSHAKE_PROTO_STATUS = 1
	HANDSHAKE_PROTO_LOGIN  = 2
)

type MinecraftProxyServer struct {
	Listen string
	Remote string
	MOTD   string

	running bool
	server  *mcnet.Listener

	HandleEncryption func(serverId string) error
	HandleLogin      func(packet *PacketLoginStart)
}

func (s *MinecraftProxyServer) StartServer() error {
	var err error
	s.server, err = mcnet.ListenMC(s.Listen)
	if err != nil {
		return err
	}
	s.running = true

	for s.running {
		conn, err := s.server.Accept()
		if err != nil {
			continue
		}
		go func() {
			s.handleConnection(&conn)
		}()
	}

	return nil
}

func (s *MinecraftProxyServer) CloseServer() {
	if s.running {
		s.server.Close()
	}
}

func (s *MinecraftProxyServer) handleConnection(conn *mcnet.Conn) error {
	defer conn.Close()
	handshake, err := ReadHandshake(conn)
	if err != nil {
		return err
	}

	if handshake.NextState == HANDSHAKE_PROTO_LOGIN {
		err = s.forwardConnection(conn, *handshake)
		return err
	} else if handshake.NextState == HANDSHAKE_PROTO_STATUS {
		err = s.handlePing(conn, *handshake)
		return err
	}

	return nil
}

// forward connection to real server
func (s *MinecraftProxyServer) forwardConnection(conn *mcnet.Conn, handshake PacketHandshake) error {
	remoteConn, err := mcnet.DialMC(s.Remote)
	if err != nil {
		return err
	}
	defer remoteConn.Close()

	// modify & send handshake packet
	if strings.Contains(handshake.ServerAddress, "\x00FML\x00") {
		handshake.ServerAddress = strings.SplitN(s.Remote, ":", 2)[0] + "\u0000FML\u0000"
	} else {
		handshake.ServerAddress = strings.SplitN(s.Remote, ":", 2)[0]
	}
	port := 25565
	{
		slice := strings.SplitN(s.Remote, ":", 2)
		if len(slice) > 1 {
			port, err = strconv.Atoi(slice[1])
			if err != nil {
				port = 25565
			}
		}
	}
	handshake.ServerPort = uint16(port)
	WriteHandshake(remoteConn, handshake)

	// read username
	loginStart, err := ReadLoginStart(conn)
	if err != nil {
		return err
	}
	if s.HandleLogin != nil {
		s.HandleLogin(loginStart)
	}
	log.Println("login:", loginStart.Name)

	WriteLoginStart(remoteConn, *loginStart)

	if s.HandleEncryption != nil {
		err = s.handleEncryption(conn, remoteConn)
		if err != nil {
			println(err.Error())
			return err
		}
	}

	log.Println("Forwarding packets")

	// forward connection
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		io.Copy(remoteConn, conn)
		wg.Done()
	}()
	go func() {
		io.Copy(conn, remoteConn)
		wg.Done()
	}()
	wg.Wait()

	return nil
}

func (s *MinecraftProxyServer) handleEncryption(conn *mcnet.Conn, remoteConn *mcnet.Conn) error {
	var p packet.Packet

	err := remoteConn.ReadPacket(&p)
	if err != nil {
		return err
	}

	if p.ID != 0x01 { // not an encryption request packet
		conn.WritePacket(p)
		return nil
	}

	pk, err := ReadEncryptionRequest(p)
	if err != nil {
		return err
	}

	key, encoStream, decoStream := newSymmetricEncryption()
	realServerId := authDigest(pk.ServerID, key, pk.PublicKey)
	err = s.HandleEncryption(realServerId)
	if err != nil {
		return err
	}

	p, err = genEncryptionKeyResponse(key, pk.PublicKey, pk.VerifyToken)
	if err != nil {
		return fmt.Errorf("gen encryption key response fail: %v", err)
	}

	err = remoteConn.WritePacket(p)
	if err != nil {
		return err
	}

	remoteConn.SetCipher(encoStream, decoStream)

	return nil
}

// handle ping request
func (s *MinecraftProxyServer) handlePing(conn *mcnet.Conn, handshake PacketHandshake) error {
	for {
		var p packet.Packet
		err := conn.ReadPacket(&p)
		if err != nil {
			return err
		}

		switch p.ID {
		case 0x00: // status request
			resp := StatusResponse{}

			resp.Version.Name = "Maid"
			resp.Version.Protocol = int(handshake.ProtocolVersion)
			resp.Players.Max = 20
			resp.Players.Online = 0
			resp.Description = s.MOTD
			if resp.Description == "" {
				resp.Description = "A Maid powered proxy server"
			}
			bytes, err := json.Marshal(resp)
			if err != nil {
				return nil
			}

			err = WriteStatusResponse(conn, PacketStatusResponse{
				Response: string(bytes),
			})
			if err != nil {
				return err
			}

		case 0x01: // ping
			var payload packet.Long
			err := p.Scan(&payload)
			if err != nil {
				return err
			}

			err = conn.WritePacket(packet.Marshal(
				0x01,
				packet.Long(payload)),
			)
			if err != nil {
				return err
			}
		}
	}

}
