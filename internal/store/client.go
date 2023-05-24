package store

import "github.com/JunNishimura/Goit/internal/file"

type Client struct {
	Conf         *Config
	Idx          *Index
	RootGoitPath string
	Head
}

func NewClient(path string) (*Client, error) {
	rootGoitPath, _ := file.FindGoitRoot(path) // ignore the error since the error is not important
	config, err := NewConfig(rootGoitPath)
	if err != nil {
		return nil, err
	}
	index, err := NewIndex(rootGoitPath)
	if err != nil {
		return nil, err
	}
	head, err := NewHead(rootGoitPath)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Conf:         config,
		Idx:          index,
		RootGoitPath: rootGoitPath,
		Head:         head,
	}

	return client, nil
}
