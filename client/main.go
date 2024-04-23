package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var port string
	var addr string
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	flag.StringVar(&port, "port", ":8000", "Porta que o servidor está rodando")
	flag.StringVar(&addr, "addr", "localhost", "Endereço em que o servidor está hospedado")
	flag.Parse()
	for {
		fmt.Println("Conectando no Servidor...")
		conn, err := net.Dial("tcp", port)
		if err != nil {
			fmt.Println("Houve um erro ao tentar estabelecer uma conexão\n\t", err)
			fmt.Println("Tentando de novo...")
			time.Sleep(time.Second * 1)
			continue
		}

		fmt.Println("Conexão aceita")
		fmt.Print("Digite seu nome público: ")
		reader := bufio.NewReader(os.Stdin)
		name, _ := reader.ReadString('\n')
		conn.Write([]byte(name))
		go sendMessage(conn)
		go toReceiveMessage(conn)
		<-signalCh
		os.Exit(0)
	}
}

func sendMessage(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	m, _ := reader.ReadString('\n')
	conn.Write([]byte(m))
	sendMessage(conn)
}

func toReceiveMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)
	m, _ := reader.ReadString('\000')
	fmt.Println(m)
	toReceiveMessage(conn)
}
