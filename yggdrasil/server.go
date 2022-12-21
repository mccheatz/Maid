package yggdrasil

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type YggdrasilHandlerJoinServer func(request YggdrasilJoinServerRequest)
type YggdrasilHandlerProfile func(uuid string) YggdrasilProfileResponse

type YggdrasilServer struct {
	Address string
	server  *http.Server

	// handlers
	HandlerJoinServer YggdrasilHandlerJoinServer
	HandlerProfile    YggdrasilHandlerProfile
}

type YggdrasilJoinServerRequest struct {
	AccessToken     string `json:"accessToken"`
	SelectedProfile string `json:"selectedProfile"`
	ServerId        string `json:"serverId"`
}

func serverBrandMiddleware(c *gin.Context) {
	c.Header("Server", "MockYggdrasil/1.0.0")
}

func (s *YggdrasilServer) StartServer() error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(serverBrandMiddleware)

	r.POST("/session/minecraft/join", func(c *gin.Context) {
		if s.HandlerJoinServer == nil {
			c.AbortWithError(http.StatusNotFound, errors.New("no handler found"))
			return
		}

		postBody, _ := ioutil.ReadAll(c.Request.Body)
		var req YggdrasilJoinServerRequest
		json.Unmarshal(postBody, &req)
		s.HandlerJoinServer(req)
	})

	r.GET("/session/minecraft/profile/:uuid", func(c *gin.Context) {
		if s.HandlerProfile == nil {
			c.AbortWithError(http.StatusNotFound, errors.New("no handler found"))
			return
		}

		resp := s.HandlerProfile(c.Param("uuid"))
		body, _ := json.Marshal(resp)

		c.String(http.StatusOK, string(body))
	})

	s.server = &http.Server{
		Addr:    s.Address,
		Handler: r,
	}

	err := s.server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}

	return err
}

func (s *YggdrasilServer) StopServer() {
	if s.server != nil {
		s.server.Close()
	}
}
