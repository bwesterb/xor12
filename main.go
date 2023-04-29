package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/cloudflare/circl/xof/k12"
	"github.com/templexxx/xor"
)

var (
	listenAddr   = flag.String("l", ":42422", "listen addr")
	upstreamAddr = flag.String("u", "", "upstream addr")
	key          = flag.String("k", "", "key")
)

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("Listen: %v", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Accept: %v", err)
		}

		go handle(conn)
	}
}

func pump(ctx string, r io.Reader, w io.Writer) {
	var buf [8192]byte
	var mask [8192]byte

	h := k12.NewDraft10([]byte{})
	_, _ = h.Write([]byte(*key))

	for {
		n, err := r.Read(buf[:])
		if err != nil {
			log.Printf("%s %v", ctx, err)
			return
		}

		_, _ = h.Read(mask[:n])

		xor.BytesSameLen(buf[:n], buf[:n], mask[:n])

		_, err = w.Write(buf[:n])
		if err != nil {
			log.Printf("%s %v", ctx, err)
			return
		}
	}
}

func handle(conn net.Conn) {
	up, err := net.Dial("tcp", *upstreamAddr)
	if err != nil {
		log.Printf("%s %v", conn.RemoteAddr(), err)
		return
	}

	go pump(fmt.Sprintf("%v up->down", conn.RemoteAddr()), up, conn)
	go pump(fmt.Sprintf("%v down->up", conn.RemoteAddr()), conn, up)
}
