package config

import (
    "encoding/json"
    "flag"
    "os"
    "path"
)

var ConfigFileName = flag.String("bacs-proto-config", "BacsProtobuf.json",
    "Name of configuration file")

type Config struct {
    Local struct {
        ProtoPrefix string
        GoPrefix    string
        Import      []struct {
            Prefix string
            Path   string
        }
    }
    Dependencies []string
}

func NewConfig() *Config {
    return &Config{}
}

func ParseConfig(config string) (cfg *Config, err error) {
    file, err := os.Open(config)
    if err != nil {
        return
    }
    defer file.Close()
    decoder := json.NewDecoder(file)
    cfg = NewConfig()
    err = decoder.Decode(cfg)
    if err != nil {
        cfg = nil
    }
    return
}

func ParseProject(project string) (*Config, error) {
    return ParseConfig(path.Join(project, *ConfigFileName))
}
