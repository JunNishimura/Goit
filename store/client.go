package store

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
	rootDir := "" // TODO: implement getRootDir function

	client := &Client{
		Conf:    config,
		Idx:     index,
		rootDir: rootDir,
	}

	return client, nil
}
