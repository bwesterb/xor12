package main

import (
	"fmt"
	"os"

	"github.com/cloudflare/circl/xof/k12"
	"github.com/templexxx/xor"
)

func main() {
	var buf [8192]byte
	var mask [8192]byte
	h := k12.NewDraft10([]byte{})

	if len(os.Args) >= 2 {
		_, _ = h.Write([]byte(os.Args[1]))
	}

	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Stdin.Read: %v", err)
			return
		}

		_, _ = h.Read(mask[:n])

		xor.BytesSameLen(buf[:n], buf[:n], mask[:n])

		_, err = os.Stdout.Write(buf[:n])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Stdout.Write: %v", err)
			return
		}
	}
}
