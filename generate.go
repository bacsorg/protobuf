package main

import (
    "gopkg.in/codegangsta/cli.v1"
    "os"
)

func main() {
    app := cli.NewApp()
    app.Name = "Google Protobuf generator"
    app.Usage = "Call from go generate to update Go protobufs"
    app.Flags = []cli.Flag{
        cli.StringFlag{Name: "config", Value: "BacsProtobuf.json"},
    }
    app.Action = func(c *cli.Context) {
        println("TODO")
    }
    app.Run(os.Args)
}
