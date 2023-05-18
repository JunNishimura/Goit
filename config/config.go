package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type KV map[string]string

type Config struct {
	Map map[string]KV
}

func NewConfig() (*Config, error) {
	config := &Config{
		Map: make(map[string]KV),
	}
	// if config file exists, read config
	configPath := filepath.Join(".goit", "config")
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		if err := config.read(); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (c *Config) read() error {
	return nil
}

func (c *Config) Add(ident, key, value string) {
	if _, ok := c.Map[ident]; ok {
		c.Map[ident][key] = value
	} else {
		c.Map[ident] = make(KV)
		c.Map[ident][key] = value
	}
}

func (c *Config) Write() error {
	configPath := filepath.Join(".goit", "config")
	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	var content string
	for ident, kv := range c.Map {
		content += fmt.Sprintf("[%s]\n", ident)
		for k, v := range kv {
			content += fmt.Sprintf("\t%s = %s\n", k, v)
		}
	}

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
