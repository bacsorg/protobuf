package main

import (
	"flag"
	"log"

	"github.com/bunsanorg/protoutils/generator"
)

func main() {
	flag.Parse()
	err := generator.Generate()
	if err != nil {
		log.Fatal(err)
	}
}
