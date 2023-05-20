package store

import "github.com/JunNishimura/Goit/file"

type Client struct {
	Conf    *Config
	Idx     *Index
	RootDir string
}

func NewClient() (*Client, error) {
	rootDir, _ := file.FindGoitRoot(".") // ignore the error since the error is not important
	config, err := NewConfig(rootDir)
	if err != nil {
		return nil, err
	}
	index, err := NewIndex()
	if err != nil {
		return nil, err
	}

	client := &Client{
		Conf:    config,
		Idx:     index,
		RootDir: rootDir,
	}

	return client, nil
}
