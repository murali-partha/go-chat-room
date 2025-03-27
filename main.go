package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	if os.Args[1] == "start" {
		server()
	}
	if os.Args[1] == "join" {
		if len(os.Args) < 3 {
			fmt.Println("Please provide a name to join the chat room")
			return
		}
		client(os.Args[2])
	} else {
		fmt.Println("Invalid mode : \"", os.Args[1], "\". Valid modes are start and join")

	}
}

func server() {
	listener, err := net.Listen("tcp", "127.0.0.1:3000")

	if err != nil {
		fmt.Println("Unable to start a TCP listener. err: ", err)
		return
	}

	defer listener.Close()

	connMap := make(map[string]net.Conn)

	for {
		conn, err := listener.Accept()

		if err != nil {

			fmt.Println("Unable to create a TCP connetion. err: ", err)
			return
		}
		buf := make([]byte, 1024)
		name, err := conn.Read(buf)
		if err != nil {

			if err == io.EOF {
				fmt.Println("Connection closed by client")
				break
			}
		}
		for connName, connItem := range connMap {

			if connName != string(buf[:name]) {
				connItem.Write([]byte("New user joined: " + string(buf[:name])))
			}
		}

		connMap[string(buf[:name])] = conn

		go func() {
			for {
				data, err := conn.Read(buf)
				if err != nil {

					if err == io.EOF {
						fmt.Println("Connection closed by client")
						break
					}
				}
				for connName, connItem := range connMap {

					if connName != string(buf[:name]) {
						connItem.Write([]byte(buf[:data]))
					}
				}
			}
		}()

	}

	for _, connItem := range connMap {

		connItem.Close()
	}
}

func client(name string) {

	scanner := bufio.NewScanner(os.Stdin)

	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println("Unable to create a TCP connection. err: ", err)
		return
	}

	go func() {

		for {
			buf := make([]byte, 1024)

			data, err := conn.Read(buf)

			if err != nil {

				if err == io.EOF {
					fmt.Println("Connection closed by client")
					break
				}
			}

			fmt.Println(string(buf[:data]))
		}
	}()

	conn.Write([]byte(name))

	for scanner.Scan() {

		conn.Write([]byte(name + ": " + scanner.Text()))
	}

	conn.Write([]byte("Bye"))

	conn.Close()
}
