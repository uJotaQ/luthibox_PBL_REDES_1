package network

import (
	"fmt"
	"net"
)

type Server struct {
	listener net.Listener
	port     string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", "0.0.0.0:"+s.port) // <- escuta em todas interfaces
	if err != nil {
		return fmt.Errorf("erro ao iniciar servidor: %v", err)
	}

	fmt.Printf("🎮 Servidor LuthiBOX rodando na porta %s...\n", s.port)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Erro ao aceitar conexão: %v\n", err)
			continue
		}

		go handleClient(conn)
	}
}
