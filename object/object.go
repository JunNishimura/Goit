package object

type ObjectType string

const (
	OBJ_DIR            = ".goit/objects"
	Blob    ObjectType = "blob"
	Tree    ObjectType = "tree"
	Commit  ObjectType = "commit"
	Tag     ObjectType = "tag"
)
