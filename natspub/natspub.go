package main

import (
	"log"
	"os"

	stan "github.com/nats-io/stan.go"
)

func main() {
	model, err := os.ReadFile("./natspub/model.json")
	if err != nil {
		log.Fatal(err)
	}

	sc, _ := stan.Connect("test-cluster", "pub")

	defer sc.Close()

	sc.Publish("Message", model)
}
