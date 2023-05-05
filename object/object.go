package object

type ObjectType int

const (
	OBJ_DIR            = ".goit/objects"
	Blob    ObjectType = iota
	Tree
	Commit
	Tag
)

func (ot ObjectType) String() string {
	switch ot {
	case Blob:
		return "blob"
	case Tree:
		return "tree"
	case Commit:
		return "commit"
	case Tag:
		return "tag"
	default:
		return ""
	}
}
