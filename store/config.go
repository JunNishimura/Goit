package store

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	identRegexp          = regexp.MustCompile(`^\[.*\]$`)
	ErrInvalidIdentifier = errors.New("fatal: invalid identifier")
)

type KV map[string]string

type Config struct {
	Map map[string]KV
}

func NewConfig() (*Config, error) {
	var config *Config
	configPath := filepath.Join(".goit", "config")
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		config, err = load()
		if err != nil {
			return nil, err
		}
	} else {
		config = newConfig()
	}
	return config, nil
}

func newConfig() *Config {
	return &Config{
		Map: make(map[string]KV),
	}
}

func load() (*Config, error) {
	config := newConfig()

	configPath := filepath.Join(".goit", "config")
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var ident string
	buf := bytes.NewReader(b)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		text := scanner.Text()
		if isIdent := identRegexp.MatchString(text); isIdent {
			if len(text) < 2 {
				return nil, ErrInvalidIdentifier
			}
			ident = text[1 : len(text)-1]
			config.Map[ident] = make(KV)
		} else {
			splitText := strings.Split(strings.Replace(text, "\t", "", -1), "=")
			key := strings.TrimSpace(splitText[0])
			value := strings.TrimSpace(splitText[1])
			config.Map[ident][key] = value
		}
	}

	return config, err
}

func (c *Config) IsUserSet() bool {
	kv, ok := c.Map["user"]
	if !ok {
		return false
	}
	if _, ok := kv["name"]; !ok {
		return false
	}
	if _, ok := kv["email"]; !ok {
		return false
	}
	return true
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
