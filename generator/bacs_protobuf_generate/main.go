package main

import (
    "flag"
    "github.com/bacsorg/protobuf/generator"
    "log"
)

func main() {
    flag.Parse()
    err := generator.Generate()
    if err != nil {
        log.Fatal(err)
    }
}
