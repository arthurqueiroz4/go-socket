package broadcast

import (
	"fmt"
	"net"
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
	sender := b.getSender(from)

	for _, recipient := range b.c {
		if recipient.conn == from {
			continue
		}
		recipient.conn.Write([]byte(fmt.Sprintf("%s:\n\t%s", sender, m)))
	}
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
