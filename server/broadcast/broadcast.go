package broadcast

import (
	"fmt"
	"net"
	"strings"
	"time"
	"unicode/utf8"
)

type Broadcast struct {
	c []Conn
}

type Conn struct {
	name string
	conn net.Conn
}

func New() *Broadcast {
	return &Broadcast{
		make([]Conn, 0),
	}
}

func (b *Broadcast) Send(m string, from net.Conn) {
	err := b.validateMessage(m)
	if err != nil {
		from.Write([]byte(err.Error()))
		return
	}
	
	sender := b.getSender(from)

	msg := fmt.Sprintf("%s às %s:\n\t%s\000", sender, time.Now().Format("02/01/06 03:04 Mon"), m)

	fmt.Println(msg)

	for _, recipient := range b.c {
		if recipient.conn == from {
			continue
		}
		recipient.conn.Write([]byte(msg))
	}
}

func (b *Broadcast) validateMessage(m string) error {
	m = strings.TrimSpace(m)
	if m == "" {
		return fmt.Errorf("a mensagem não pode estar vazia\000")
	}

	if utf8.RuneCountInString(m) > 100 {
		return fmt.Errorf("a mensagem não pode ter mais do que 100 caracteres\000")
	}

	return nil
}

func (b *Broadcast) getSender(from net.Conn) string {
	for _, recipient := range b.c {
		if recipient.conn == from {
			return recipient.name
		}
	}
	return ""
}

func (b *Broadcast) Add(n string, c net.Conn) {
	b.c = append(b.c, Conn{n, c})
}

func (b *Broadcast) Remove(toRemove net.Conn) {
	for i, c := range b.c {
		if c.conn == toRemove {
			fmt.Println("Encerrando conexão com", c.name)
			b.removeByIndex(i)
		}
	}
}

func (b *Broadcast) removeByIndex(i int) {
	b.c = append(b.c[:i], b.c[i+1:]...)
}
