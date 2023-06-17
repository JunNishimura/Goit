package store

type Client struct {
	Conf         *Config
	Idx          *Index
	Head         *Head
	Refs         *Refs
	Ignore       *Ignore
	RootGoitPath string
}

func NewClient(config *Config, index *Index, head *Head, refs *Refs, ignore *Ignore, rootGoitPath string) *Client {
	return &Client{
		Conf:         config,
		Idx:          index,
		Head:         head,
		Refs:         refs,
		Ignore:       ignore,
		RootGoitPath: rootGoitPath,
	}
}
