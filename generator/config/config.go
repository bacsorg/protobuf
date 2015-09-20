package config

import (
    "encoding/json"
    "os"
)

type Config struct {
    LocalImportPrefix string
}

func NewConfig() *Config {
    return &Config{}
}

func ParseConfig(path string) (cfg *Config, err error) {
    file, err := os.Open(path)
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
