package config

type Config struct {
    LocalImportPrefix string
}

func NewConfig() *Config {
    return &Config{}
}

func ParseConfig(path string) (cfg *Config, err error) {
    cfg = NewConfig()
    return
}
