package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	remoteServer, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	fmt.Println("LISTENING 8081")

	localServer, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("LISTENING 8080")

	conn, err := remoteServer.Accept()
	if err != nil {
		panic(err)
	}

	for {
		fowardConn, err := createFowardConn(conn)
		if err != nil {
			panic(err)
		}

		handle(fowardConn, localServer)
	}
}

func createFowardConn(conn net.Conn) (net.Conn, error) {
	server, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalln(err)
	}

	port := server.Addr().(*net.TCPAddr).Port
	header := fmt.Sprintf("handshake-port:%d\n", port)
	conn.Write([]byte(header))

	fowardConn, err := server.Accept()
	if err != nil {
		return nil, err
	}

	return fowardConn, nil
}

func handle(remote net.Conn, server net.Listener) {
	localConn, err := server.Accept()
	if err != nil {
		panic(err)
	}

	go func() {
		_, err = io.Copy(remote, localConn)
		if err != nil {
			panic(err)
		}

		localConn.Close()
	}()

	go func() {
		_, err = io.Copy(localConn, remote)
		if err != nil {
			panic(err)
		}

		remote.Close()
	}()
}
