package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/arthurqueiroz04/go-socket/broadcast"
)

func getNameFromConn(c net.Conn) string {
	msg := "Digite seu nome: "
	c.Write([]byte(msg))
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
	fmt.Println("Escutando a porta 8000")

	ln, _ := net.Listen("tcp", ":8000")

	bc := broadcast.New()

	for {
		c, _ := ln.Accept()
		name := getNameFromConn(c)
		fmt.Println("Conexão aceita com", name)
		bc.Add(name, c)
		go handler(c, bc)
	}

}
