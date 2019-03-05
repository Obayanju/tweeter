package main

import "log"

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
	state := new(Feed)

	var junk Nothing
	if err := state.Post("Hello world!", &junk); err != nil {
		log.Fatalf("Post: %v", err)
	}
	if err := state.Post("Today is Monday", &junk); err != nil {
		log.Fatalf("Post: %v", err)
	}

	var lst []string
	if err := state.Get(5, &lst); err != nil {
		log.Fatalf("Get: %v", err)
	}

	for _, elt := range lst {
		log.Println(elt)
	}
}
