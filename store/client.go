package store

type Client struct {
	*Config
	*Index
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
	rootDir := "" // TODO: implement getRootDir function

	client := &Client{
		Config:  config,
		Index:   index,
		rootDir: rootDir,
	}

	return client, nil
}
