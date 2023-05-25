package store

type Client struct {
	Conf *Config
	Idx  *Index
	Head
	RootGoitPath string
}

func NewClient(config *Config, index *Index, head Head, rootGoitPath string) *Client {
	return &Client{
		Conf:         config,
		Idx:          index,
		Head:         head,
		RootGoitPath: rootGoitPath,
	}
}
