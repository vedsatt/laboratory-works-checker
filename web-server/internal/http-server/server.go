package server

import "web_server/internal/config"

type Server struct {
	port string
	url  string
}

func New() *Server {
	cfg := config.GetConfig()

	s := &Server{
		port: cfg.Port,
		url:  cfg.URL,
	}
	return s
}

func (s *Server) Start() {

}
