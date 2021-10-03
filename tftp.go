package main

import (
	"flag"
	"io/ioutil"
	"log"
	tftp "network_go/ch6/tftp"
)

var (
	address = flag.String("a", "127.0.0.1:6667", "listen address")
	payload = flag.String("p", "payload.txt", "file to serve to clients")
)

func main() {
	flag.Parse()

	p, err := ioutil.ReadFile(*payload)
	if err != nil {
		log.Fatal(err)
	}

	s := tftp.Server{Payload: p}
	log.Fatal(s.ListenAndServe(*address))
}
