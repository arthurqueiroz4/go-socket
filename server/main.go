package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/arthurqueiroz04/go-socket/server/broadcast"
)

func getNameFromConn(c net.Conn) string {
	name, _ := bufio.NewReader(c).ReadString('\n')
	return strings.TrimSpace(name)
}

func handler(conn net.Conn, bc *broadcast.Broadcast) {
	for {
		m, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				bc.Remove(conn)
				conn.Close()
				return
			}
			fmt.Println("Error reading from connection", err)
			return
		}
		bc.Send(m, conn)
	}
}

func main() {
	var port string
	flag.StringVar(&port, "port", ":8000", "Porta que será escutada pelo programa")
	flag.Parse()

	fmt.Println("Escutando a porta", port)

	ln, _ := net.Listen("tcp", port)

	bc := broadcast.New()

	for {
		c, _ := ln.Accept()
		name := getNameFromConn(c)
		fmt.Println("Conexão aceita com", name)
		bc.Add(name, c)
		go handler(c, bc)
	}

}
