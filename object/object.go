package object

type ObjectType string

var (
	Blob   ObjectType = "blob"
	Tree   ObjectType = "tree"
	Commit ObjectType = "commit"
	Tag    ObjectType = "tag"
)
