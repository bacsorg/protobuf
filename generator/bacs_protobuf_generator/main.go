package main

import (
    "github.com/bacsorg/protobuf/generator"
    "github.com/bacsorg/protobuf/generator/config"
    "gopkg.in/codegangsta/cli.v1"
    "log"
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
        cfg, err := config.ParseConfig(c.String("config"))
        if err != nil {
            log.Fatal(err)
        }
        err = generator.Generate(cfg)
        if err != nil {
            log.Fatal(err)
        }
    }
    app.Run(os.Args)
}
