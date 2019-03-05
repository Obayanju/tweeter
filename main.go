package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"flag"
	"os"
)

const (
	defaultHost = "localhost"
	defaultPort = "3410"
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
	var address string
	var isClient bool
	var isServer bool
	flag.BoolVar(&isClient, "client", false, "starts as tweeter client")
	flag.BoolVar(&isServer, "server", false, "starts as tweeter server")
	flag.Parse()

	if isServer && isClient {
		printUsage()
	}
	if !isServer && !isClient {
		printUsage()
	}

	switch flag.NArg() {
	case 0:
		if isClient {
			address = defaultHost + ":" + defaultPort
		} else {
			address = ":" + defaultPort
		}

	case 1:
		address = flag.Arg(0)

	default:
		printUsage()
	}

	if isClient {
		client(address)
	} else {
		server(address)
	}
}

func printUsage() {
	log.Printf("Usage: %s [-server or -client] [address]", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func server(address string) {
	// create instance of object to be exported 
	feed := new(Feed)
	// register with RPC library
	rpc.Register(feed)
	// handle http request once server is up and running
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", address)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	if err := http.Serve(l, nil); err != nil {
		log.Fatalf("http.Serve: %v", err)
	}
}

func client(address string) {
	var junk Nothing
	if err := call(address, "Feed.Post", "Hi there", &junk); err != nil {
		log.Fatalf("client.Post: %v", err)
	}
	if err := call(address, "Feed.Post", "RPC is so fun", &junk); err != nil {
		log.Fatalf("client.Post: %v", err)
	}
	var replyList []string
	if err := call(address, "Feed.Get", 4, &replyList); err != nil {
		log.Fatalf("client.Get: %v", err)
	}
	for _, elt := range replyList {
		log.Println(elt)
	}
}

func call(address string, method string, request interface{}, response interface{}) error {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Printf("rpc.DialHTTP: %v", err)
		return err
	}
	defer client.Close()

	if err = client.Call(method, request, response); err != nil {
		log.Printf("client.Call %s: %v", method, err)
		return err
	}

	return nil
}
