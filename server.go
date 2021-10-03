package tftp

import (
	"errors"
	"log"
	"net"
	"time"
)

type Server struct {
	Payload []byte        // payload that will be served
	Retries uint8         // how many times to retry failed transmissions
	Timeout time.Duration // how long to wait for an acknowledgement
}

func (s Server) ListenAndServe(addr string) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	log.Printf("Listening on %s \n")

	return s.Serve(conn)
}

func (s *Server) Serve(conn net.PacketConn) error {
	if conn == nil {
		return errors.New("No Connection")
	}

	if s.Payload == nil {
		return errors.New("payload is empty")
	}

	if s.Retries == 0 {
		s.Retries = 10
	}

	if s.Timeout == 0 {
		s.Timeout = 6 * time.Second
	}

	var rrq ReadReq

	for {
		buf := make([]byte, DatagramSize)

		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return nil
		}

		err = rrq.UnmarshalBinary(buf)
		if err != nil {
			log.Printf("[%s] bad request: %v", addr, err)
			continue
		}

		go s.handle(addr.String(), rrq)
	}
}

func (s Server) handle(clientAddr string, rrq ReadReq) {
	log.Printf("[%s] file requested: %s", clientAddr, rrq.Filename)
}
