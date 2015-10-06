package main

import (
    "flag"
    "log"

    "github.com/bacsorg/protobuf/generator"
)

func main() {
    flag.Parse()
    err := generator.Generate()
    if err != nil {
        log.Fatal(err)
    }
}
