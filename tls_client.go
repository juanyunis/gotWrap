package gotWrap

import (
    "crypto/tls"
    "io"
    "log"
)

type Client struct {
    RemoteAddr string
    Protocol string
    PemFile string
    KeyFile string
}

func (client *Client) Connect() {
    cert, err := tls.LoadX509KeyPair(client.PemFile, client.KeyFile)
    if err != nil {
        log.Fatalf("server: loadkeys: %s", err)
    }
    config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
    conn, err := tls.Dial(client.Protocol, client.RemoteAddr, &config)
    if err != nil {
        log.Fatalf("client: dial: %s", err)
    }
    defer conn.Close()
    log.Println("client: connected to: ", conn.RemoteAddr())
    state := conn.ConnectionState()
    log.Println("client: handshake: ", state.HandshakeComplete)
    log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
    message := "Hello\n"
    n, err := io.WriteString(conn, message)
    if err != nil {
        log.Fatalf("client: write: %s", err)
    }
    log.Printf("client: wrote %q (%d bytes)", message, n)
    reply := make([]byte, 256)
    n, err = conn.Read(reply)
    log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
    log.Print("client: exiting")
}
