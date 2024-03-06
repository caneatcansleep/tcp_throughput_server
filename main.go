package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	// Create server
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panicln(err)
	}

	go func() {
		<-stopper
		log.Println("Stoping..")
		os.Exit(0)
	}()

	log.Println("Listening..")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("New connection", conn.RemoteAddr())
		go handleEcho(conn)
	}
}

func handleEcho(conn net.Conn) {
	var cnt uint64 = 0
	buf := make([]byte, 256*1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		cnt += uint64(n)
	}
	// log.Println(cnt)
}
