package store

type Client struct {
	Conf         *Config
	Idx          *Index
	Head         *Head
	Refs         *Refs
	RootGoitPath string
}

func NewClient(config *Config, index *Index, head *Head, refs *Refs, rootGoitPath string) *Client {
	return &Client{
		Conf:         config,
		Idx:          index,
		Head:         head,
		Refs:         refs,
		RootGoitPath: rootGoitPath,
	}
}
