package object

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/JunNishimura/Goit/internal/binary"
	"github.com/JunNishimura/Goit/internal/sha"
)

type Node struct {
	Hash     sha.SHA1
	Name     string
	Children []*Node
}

func (n *Node) GetPaths() []string {
	if len(n.Children) == 0 {
		return []string{n.Name}
	}

	var paths []string
	for _, child := range n.Children {
		ps := child.getPaths(n.Name)
		paths = append(paths, ps...)
	}

	return paths
}

func (n *Node) getPaths(parentDir string) []string {
	if len(n.Children) == 0 {
		return []string{fmt.Sprintf("%s/%s", parentDir, n.Name)}
	}

	var paths []string
	for _, child := range n.Children {
		ps := child.getPaths(fmt.Sprintf("%s/%s", parentDir, n.Name))
		paths = append(paths, ps...)
	}

	return paths
}

type Tree struct {
	object   *Object
	Children []*Node
}

func NewTree(rootGoitPath string, object *Object) (*Tree, error) {
	if object.Type != TreeObject {
		return nil, ErrInvalidTreeObject
	}
	t := newTree(object)
	if err := t.load(rootGoitPath); err != nil {
		return nil, ErrInvalidTreeObject
	}
	return t, nil
}

func newTree(object *Object) *Tree {
	return &Tree{
		object:   object,
		Children: []*Node{},
	}
}

func (t *Tree) load(rootGoitPath string) error {
	children, err := walkTree(rootGoitPath, t.object)
	if err != nil {
		return err
	}
	t.Children = children
	return nil
}

func walkTree(rootGoitPath string, object *Object) ([]*Node, error) {
	var nodes []*Node
	var isDir bool
	var nodeName string
	isFirstLine := true

	buf := bytes.NewReader(object.Data)
	for {
		var lineSplit []string
		if isFirstLine {
			lineString, err := binary.ReadNullTerminatedString(buf)
			if err != nil {
				return nil, err
			}
			lineSplit = strings.Split(lineString, " ")

			mode := lineSplit[0]
			if mode == "040000" {
				isDir = true
			}
			nodeName = lineSplit[1]

			isFirstLine = false
		} else {
			// get 20 bytes to read hash
			hashBytes := make([]byte, 20)
			n, err := buf.Read(hashBytes)
			if err != nil {
				return nil, err
			}
			if n != 20 {
				return nil, errors.New("fail to read hash")
			}

			// read filemode and filename
			lineString, err := binary.ReadNullTerminatedString(buf)
			if err != nil {
				return nil, err
			}

			// append lineSplit
			hashString := hex.EncodeToString(hashBytes)
			lineSplit = []string{hashString}
			if lineString != "" {
				lineSplit = append(lineSplit, strings.Split(lineString, " ")...)
			}

			hash, err := sha.ReadHash(hashString)
			if err != nil {
				return nil, err
			}
			var children []*Node
			if isDir {
				treeObject, err := GetObject(rootGoitPath, hash)
				if err != nil {
					return nil, err
				}
				if treeObject.Type != TreeObject {
					return nil, ErrInvalidTreeObject
				}
				gotChildren, err := walkTree(rootGoitPath, treeObject)
				if err != nil {
					return nil, err
				}
				children = gotChildren
				isDir = false
			}
			node := &Node{
				Hash:     hash,
				Name:     nodeName,
				Children: children,
			}
			nodes = append(nodes, node)

			// last line
			if len(lineSplit) == 1 {
				break
			}

			mode := lineSplit[1]
			if mode == "040000" {
				isDir = true
			}
			nodeName = lineSplit[2]
		}
	}

	return nodes, nil
}

func (t *Tree) String() string {
	var lines []string

	for _, childNode := range t.Children {
		var line string
		if len(childNode.Children) == 0 {
			line = fmt.Sprintf("100644 blob %s\t%s", childNode.Hash, childNode.Name)
		} else {
			line = fmt.Sprintf("040000 tree %s\t%s", childNode.Hash, childNode.Name)
		}
		lines = append(lines, line)
	}

	message := strings.Join(lines, "\n")
	return message
}

func GetNode(children []*Node, path string) (*Node, bool) {
	if len(children) == 0 {
		return nil, false
	}

	pathSplit := strings.SplitN(path, "/", 2)
	searchName := pathSplit[0]

	left := 0
	right := len(children)
	for {
		middle := (left + right) / 2
		node := children[middle]
		if node.Name == searchName {
			if len(node.Children) == 0 {
				return node, true
			}
			if len(pathSplit) > 1 {
				return GetNode(node.Children, pathSplit[1])
			} else {
				return node, true
			}
		} else if node.Name < searchName {
			left = middle + 1
		} else {
			right = middle
		}

		if right-left < 1 {
			break
		}
	}

	return nil, false
}
