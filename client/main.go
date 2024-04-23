package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var errCh chan error

func main() {
	var port string
	var host string
	var maxTentatives = 5
	var count int

	signalCh := make(chan os.Signal, 1)
	errCh = make(chan error)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		os.Exit(1)
	}()

	flag.StringVar(&port, "port", "8000", "Porta que o servidor está rodando")
	flag.StringVar(&host, "host", "localhost", "Endereço em que o servidor está hospedado")
	flag.Parse()

	fmt.Println(host, port)
	for {
		if count >= maxTentatives {
			fmt.Println("Chegou ao limite de tentativas de conexão.")
			os.Exit(1)
		}
		count++
		
		fmt.Println("Conectando no Servidor...")
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
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
		err = <- errCh
		if err == io.EOF {
			fmt.Println("Conexão fechada com o servidor.")
			os.Exit(1)
		}
		fmt.Println("Ocorreu um erro:\n\t", err)
		os.Exit(2)
	}
}

func sendMessage(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	m, _ := reader.ReadString('\n')
	_, err := conn.Write([]byte(m))
	if err != nil {
		errCh <- err
	}
	sendMessage(conn)
}

func toReceiveMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)
	m, err := reader.ReadString('\000')
	if err != nil {
		errCh <- err
	}
	fmt.Println(m)
	toReceiveMessage(conn)
}
