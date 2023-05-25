package object

type Type int

const (
	UndefinedObject Type = iota
	BlobObject
	TreeObject
	CommitObject
	TagObject
)

func (t Type) String() string {
	switch t {
	case BlobObject:
		return "blob"
	case TreeObject:
		return "tree"
	case CommitObject:
		return "commit"
	case TagObject:
		return "tag"
	default:
		return "undefined"
	}
}

func NewType(typeString string) (Type, error) {
	switch typeString {
	case "blob":
		return BlobObject, nil
	case "tree":
		return TreeObject, nil
	case "commit":
		return CommitObject, nil
	case "tag":
		return TagObject, nil
	default:
		return UndefinedObject, ErrInvalidObject
	}
}
