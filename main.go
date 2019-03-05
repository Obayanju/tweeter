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
	client()
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

func client() {
	client, err := rpc.DialHTTP("tcp", "localhost"  + ":8080")
	if err != nil {
		log.Fatalf("rpc.DialHTTP: %v", err)
	}
	var junk Nothing
	if err = client.Call("Feed.Post", "Hi there", &junk); err != nil {
		log.Fatalf("client.Post: %v", err)
	}
	if err = client.Call("Feed.Post", "RPC is so fun", &junk); err != nil {
		log.Fatalf("client.Post: %v", err)
	}
	var replyList []string
	if err = client.Call("Feed.Get", 4, &replyList); err != nil {
		log.Fatalf("client.Get: %v", err)
	}
	for _, elt := range replyList {
		log.Println(elt)
	}
	if err := client.Close(); err != nil {
		log.Fatalf("client.Close: %v", err)
	}
}








