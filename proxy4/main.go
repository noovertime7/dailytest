package main

import (
	"flag"
	"log"
	"net"
)

func main() {
	help := flag.Bool("help", false, "print usage")
	bind := flag.String("bind", "127.0.0.1:6001", "The address to bind to")
	backend := flag.String("backend", "tjyw-k8s-api:6443", "The backend server address")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *backend == "" {
		flag.Usage()
		return
	}

	if *bind == "" {
		//use default bind
		log.Print("use default bind")
	}

	success, err := RunProxy(*bind, *backend)
	if !success {
		log.Fatal(err)
	}
}

func RunProxy(bind, backend string) (bool, error) {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		return false, err
	}
	defer listener.Close()
	log.Print("tcp-proxy started.")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		} else {
			go ConnectionHandler(conn, backend)
		}
	}
}

func ConnectionHandler(conn net.Conn, backend string) {
	log.Print("client connected.")
	target, err := net.Dial("tcp", backend)
	defer conn.Close()
	if err != nil {
		log.Print(err)
	} else {
		defer target.Close()
		log.Print("backend connected.")

		closed := make(chan bool, 2)
		go Proxy(conn, target, closed)
		go Proxy(target, conn, closed)
		<-closed
		log.Print("Connection closed.")
	}
}

func Proxy(from net.Conn, to net.Conn, closed chan bool) {
	buffer := make([]byte, 4096)
	for {
		n1, err := from.Read(buffer)
		if err != nil {
			closed <- true
			return
		}
		n2, err := to.Write(buffer[:n1])
		if err != nil {
			closed <- true
			return
		}
		log.Printf("n2 : %v", n2)
	}
}
