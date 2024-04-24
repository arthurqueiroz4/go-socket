package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var errCh chan error

func main() {
	var port string
	var host string
	var maxAttemps int
	errCh = make(chan error)

	flag.StringVar(&port, "port", "8000", "Porta que o servidor está rodando")
	flag.StringVar(&host, "host", "localhost", "Endereço em que o servidor está hospedado")
	flag.IntVar(&maxAttemps, "maxAttemps", 5, "Número máximo de tentativas de conexão")
	flag.Parse()

	conn, err := fistConnection(host, port, maxAttemps)
	if err != nil {
		fmt.Println("Erro ao estabelecer conexão. Err:", err)
		os.Exit(1)
	}

	greeting(conn)

	go sendMessage(conn)
	go toReceiveMessage(conn)
	handleErrors(<-errCh)
	
}

func fistConnection(host, port string, maxAttemps int) (net.Conn, error) {
	for i := 0; i < maxAttemps; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			fmt.Println("Houve um erro ao tentar estabelecer uma conexão\n\t", err)
			fmt.Println("Tentando de novo...")
			time.Sleep(time.Second * 1)
			continue
		}
		return conn, nil
	}

	return nil, fmt.Errorf("não foi possível estabelecer uma conexão com o servidor")
}

func handleErrors(err error) {
	if err == io.EOF {
		fmt.Println("Conexão fechada com o servidor.")
		os.Exit(1)
	}
	fmt.Println("Ocorreu um erro:\n\t", err)
	os.Exit(2)
}

func greeting(conn net.Conn) {
	fmt.Println("Conexão aceita")
	fmt.Print("Digite seu nome público: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	conn.Write([]byte(name))
	fmt.Println("Conectado com sucesso!")
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
