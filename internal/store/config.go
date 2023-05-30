package store

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	identRegexp          = regexp.MustCompile(`^\[.*\]$`)
	ErrInvalidIdentifier = errors.New("fatal: invalid identifier")
)

type kv map[string]string

type Config struct {
	local  map[string]kv
	global map[string]kv
}

func NewConfig(rootGoitDir string) (*Config, error) {
	config := newConfig()

	// load from global config
	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	globalConfigPath := filepath.Join(userHomePath, ".goitconfig")
	if _, err := os.Stat(globalConfigPath); !os.IsNotExist(err) {
		if err := config.load(globalConfigPath, true); err != nil {
			return nil, err
		}
	}

	// load from local config
	localConfigPath := filepath.Join(rootGoitDir, "config")
	if _, err := os.Stat(localConfigPath); !os.IsNotExist(err) {
		if err := config.load(localConfigPath, false); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func newConfig() *Config {
	return &Config{
		local:  make(map[string]kv),
		global: make(map[string]kv),
	}
}

func (c *Config) load(configPath string, isGlobal bool) error {
	b, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var ident string
	buf := bytes.NewReader(b)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		text := scanner.Text()
		if identRegexp.MatchString(text) {
			if len(text) <= 2 {
				return ErrInvalidIdentifier
			}
			ident = text[1 : len(text)-1]
			if isGlobal {
				c.global[ident] = make(kv)
			} else {
				c.local[ident] = make(kv)
			}
		} else {
			splitText := strings.Split(strings.Replace(text, "\t", "", -1), "=")
			key := strings.TrimSpace(splitText[0])
			value := strings.TrimSpace(splitText[1])
			if isGlobal {
				c.global[ident][key] = value
			} else {
				c.local[ident][key] = value
			}
		}
	}

	return nil
}

func (c *Config) IsUserSet() bool {
	localKV, localOK := c.local["user"]
	globalKV, globalOK := c.global["user"]
	if !localOK && !globalOK {
		return false
	}
	if _, ok := localKV["name"]; !ok {
		if _, ok := globalKV["name"]; !ok {
			return false
		}
	}
	if _, ok := localKV["email"]; !ok {
		if _, ok := globalKV["email"]; !ok {
			return false
		}
	}
	return true
}

func (c *Config) GetUserName() string {
	// search local config first
	localKV, localOk := c.local["user"]
	if localOk {
		if v, ok := localKV["name"]; ok {
			return v
		}
	}
	globalKV, globaOK := c.global["user"]
	if globaOK {
		if v, ok := globalKV["name"]; ok {
			return v
		}
	}
	return ""
}

func (c *Config) GetEmail() string {
	// search local config first
	localKV, localOk := c.local["user"]
	if localOk {
		if v, ok := localKV["email"]; ok {
			return v
		}
	}
	globalKV, globaOK := c.global["user"]
	if globaOK {
		if v, ok := globalKV["email"]; ok {
			return v
		}
	}
	return ""
}

func (c *Config) Add(ident, key, value string, isGlobal bool) {
	if isGlobal {
		if _, ok := c.global[ident]; ok {
			c.global[ident][key] = value
		} else {
			c.global[ident] = make(kv)
			c.global[ident][key] = value
		}
	} else {
		if _, ok := c.local[ident]; ok {
			c.local[ident][key] = value
		} else {
			c.local[ident] = make(kv)
			c.local[ident][key] = value
		}
	}
}

func (c *Config) Write(configPath string, isGlobal bool) error {
	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	var content string
	var kvs map[string]kv
	if isGlobal {
		kvs = c.global
	} else {
		kvs = c.local
	}
	for ident, keyValue := range kvs {
		content += fmt.Sprintf("[%s]\n", ident)
		for k, v := range keyValue {
			content += fmt.Sprintf("\t%s = %s\n", k, v)
		}
	}

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
