package gotWrap

import (
	"net"
	"crypto/tls"
	"log"
)

type Server struct {
	ListenerAddr string
	Protocol string
	PemFile string
	KeyFile string
}

func (server *Server) CreateServer() {
	//TODO - Auto gen certs upon first start
	cert, err := tls.LoadX509KeyPair(server.PemFile, server.KeyFile)
	if err != nil {
	log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
	listener, err := tls.Listen(server.Protocol, server.ListenerAddr, &config)
	if err != nil {
		log.Fatalf("server: listening on: %s :%s", listener.Addr().String(), err)
	}
	log.Print("server: listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	tlscon, ok := conn.(*tls.Conn)
	if ok {
		log.Print("server: conn: type assert to TLS succeedded")
		err := tlscon.Handshake()
		if err != nil {
			log.Fatalf("server: handshake failed: %s", err)
		} else {
			log.Print("server: conn: Handshake completed")
		}
		state := tlscon.ConnectionState()
		log.Println("server: mutual: ", state.NegotiatedProtocolIsMutual)
		buf := make([]byte, 512)
		for {
			log.Print("server: conn: waiting")
			n, err := conn.Read(buf)
			if err != nil {
				if err != nil {
					log.Printf("server: conn: read: %s", err)
				}
				break
 			}
			log.Printf("server: conn: echo %q\n", string(buf[:n]))
			n, err = conn.Write(buf[:n])
			log.Printf("server: conn: wrote %d bytes", n)
			if err != nil {
				log.Printf("server: write: %s", err)
				break
			}
		}
	}
	log.Println("server: conn: closed")
}