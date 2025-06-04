package server

import (
	"fmt"
	"net/http"

	libhttp "github.com/eviltomorrow/king/lib/http"
	"github.com/eviltomorrow/king/lib/http/middleware"
	"github.com/eviltomorrow/king/lib/log"
	"github.com/eviltomorrow/king/lib/network"
	"github.com/labstack/echo/v4"
)

type HTTP struct {
	network *network.Config
	log     *log.Config

	server  *http.Server
	handler *echo.Echo

	RegisteredAPI []func(libhttp.Router) error
}

func NewHTTP(network *network.Config, log *log.Config, supported ...func(libhttp.Router) error) *HTTP {
	return &HTTP{
		network: network,
		log:     log,

		RegisteredAPI: supported,

		handler: echo.New(),
	}
}

func (h *HTTP) Serve() error {
	h.handler.Use(middleware.ServerRecoveryInterceptor())
	h.handler.Use(middleware.ServerLogInterceptor())

	for _, api := range h.RegisteredAPI {
		api(h.handler)
	}

	h.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", h.network.BindIP, h.network.BindPort),
		Handler: h.handler,
	}

	return h.server.ListenAndServe()
}

func (h *HTTP) Stop() error {
	if h.server != nil {
		return h.server.Close()
	}
	return nil
}
