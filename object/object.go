package object

type ObjectType int

const (
	OBJ_DIR              = ".goit/objects"
	BLOB_TYPE ObjectType = iota
	TREE_TYPE
	COMMIT_TYPE
	TAG_TYPE
)

func (ot ObjectType) String() string {
	switch ot {
	case BLOB_TYPE:
		return "blob"
	case TREE_TYPE:
		return "tree"
	case COMMIT_TYPE:
		return "commit"
	case TAG_TYPE:
		return "tag"
	default:
		return ""
	}
}
