package store

import "github.com/JunNishimura/Goit/file"

type Client struct {
	Conf         *Config
	Idx          *Index
	RootGoitPath string
}

func NewClient() (*Client, error) {
	rootGoitPath, _ := file.FindGoitRoot(".") // ignore the error since the error is not important
	config, err := NewConfig(rootGoitPath)
	if err != nil {
		return nil, err
	}
	index, err := NewIndex(rootGoitPath)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Conf:         config,
		Idx:          index,
		RootGoitPath: rootGoitPath,
	}

	return client, nil
}
