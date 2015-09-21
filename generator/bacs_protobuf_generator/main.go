package main

import (
    "github.com/bacsorg/protobuf/generator"
    "gopkg.in/codegangsta/cli.v1"
    "log"
    "os"
)

func main() {
    app := cli.NewApp()
    app.Name = "Google Protobuf generator"
    app.Usage = "Call from go generate to update Go protobufs"
    app.Action = func(c *cli.Context) {
        err := generator.Generate()
        if err != nil {
            log.Fatal(err)
        }
    }
    app.Run(os.Args)
}
