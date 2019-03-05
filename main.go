package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// server object to be exported over RPC
type Feed struct {
	Messages []string
}

type Nothing struct{}

// methods to be exported by RPC
func (f *Feed) Post(msg string, reply *Nothing) error {
	f.Messages = append(f.Messages, msg)
	return nil
}

func (f *Feed) Get(count int, reply *[]string) error {
	if len(f.Messages) < count {
		count = len(f.Messages)
	}
	*reply = make([]string, count)
	copy(*reply, f.Messages[len(f.Messages) - count:])
	return nil
}

func main() {
	server()
}

func server() {
	// create instance of object to be exported 
	feed := new(Feed)
	// register with RPC library
	rpc.Register(feed)
	// handle http request once server is up and running
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":8080")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	if err := http.Serve(l, nil); err != nil {
		log.Fatalf("http.Serve: %v", err)
	}
}










