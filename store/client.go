package store

import "github.com/JunNishimura/Goit/file"

type Client struct {
	Conf    *Config
	Idx     *Index
	rootDir string
}

func NewClient() (*Client, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}
	index, err := NewIndex()
	if err != nil {
		return nil, err
	}
	rootDir, _ := file.FindGoitRoot(".") // ignore the error since the error is not important

	client := &Client{
		Conf:    config,
		Idx:     index,
		rootDir: rootDir,
	}

	return client, nil
}
