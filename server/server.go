package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/go-sftp-server/handler"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Serve starts an SFTP server at the given address and port using the provided SSH configuration.
func Serve(sshConfig *ssh.ServerConfig, address, port string) {
	listenAt := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp", listenAt)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", listenAt, err)
	}
	defer listener.Close()

	log.Printf("Listening on %v\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn, sshConfig)
	}
}

// handleConnection handles incoming SSH connections in parallel.
func handleConnection(conn net.Conn, sshConfig *ssh.ServerConfig) {
	defer conn.Close()

	sconn, chans, reqs, err := ssh.NewServerConn(conn, sshConfig)
	if err != nil {
		log.Printf("failed to perform SSH handshake: %v", err)
		return
	}
	defer sconn.Close()

	log.Printf("User: %s logged in.", sconn.User())
	defer log.Printf("User: %s logged out.", sconn.User())

	// Discard incoming requests to keep the connection alive
	go ssh.DiscardRequests(reqs)

	// Handle incoming channels for the session
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("User: %s Failed to accept channel: %v", sconn.User(), err)
			continue
		}

		go handleRequests(requests)
		go handleSFTP(channel, sconn.User())
	}
}

// handleRequests handles SSH session requests, specifically the "subsystem" request for SFTP.
func handleRequests(requests <-chan *ssh.Request) {
	for req := range requests {
		if req.Type == "subsystem" && string(req.Payload[4:]) == "sftp" {
			req.Reply(true, nil)
		} else {
			req.Reply(false, nil)
		}
	}
}

// handleSFTP handles the SFTP subsystem for an accepted SSH channel.
func handleSFTP(channel ssh.Channel, user string) {
	defer channel.Close()

	ctx := context.Background()
	sftpHandler, err := handler.NewSftpHandler(ctx)
	if err != nil {
		log.Fatalf("failed to initialize handler: %s", err)
	}

	server := sftp.NewRequestServer(channel, *sftpHandler)
	if err := server.Serve(); err == io.EOF {
		log.Printf("User: %s sftp client exited session.", user)
	} else if err != nil {
		log.Printf("sftp server completed with error: %s", err)
	}
}
