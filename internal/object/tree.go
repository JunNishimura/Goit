package object

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/JunNishimura/Goit/internal/binary"
	"github.com/JunNishimura/Goit/internal/sha"
	index "github.com/JunNishimura/Goit/internal/store"
)

type node struct {
	hash     sha.SHA1
	name     string
	children []*node
}

type tree struct {
	object   *Object
	children []*node
}

func NewTree(rootGoitPath string, object *Object) (*tree, error) {
	if object.Type != TreeObject {
		return nil, ErrInvalidTreeObject
	}
	t := newTree(object)
	if err := t.load(rootGoitPath); err != nil {
		return nil, ErrInvalidTreeObject
	}
	return t, nil
}

func newTree(object *Object) *tree {
	return &tree{
		object:   object,
		children: []*node{},
	}
}

func (t *tree) load(rootGoitPath string) error {
	children, err := walkTree(rootGoitPath, t.object)
	if err != nil {
		return err
	}
	t.children = children
	return nil
}

func walkTree(rootGoitPath string, object *Object) ([]*node, error) {
	var nodes []*node
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
			var children []*node
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
			node := &node{
				hash:     hash,
				name:     nodeName,
				children: children,
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

func (to *Object) ExtractEntries(rootGoitPath, rootDir string) ([]*index.Entry, error) {
	var entries []*index.Entry
	var dirName string
	var filePath string
	isFirstLine := true

	buf := bytes.NewReader(to.Data)
	for {
		var lineSplit []string
		if isFirstLine {
			lineString, err := binary.ReadNullTerminatedString(buf)
			if err != nil {
				return nil, err
			}
			lineSplit = strings.Split(lineString, " ")
			isFirstLine = false
		} else {
			// read 20 bytes sha1 hash
			hashBytes := make([]byte, 20)
			n, err := buf.Read(hashBytes)
			if err != nil {
				return nil, err
			}
			if n != 20 {
				return nil, errors.New("fail to read hash")
			}

			// read filemode and path
			lineString, err := binary.ReadNullTerminatedString(buf)
			if err != nil {
				return nil, err
			}

			lineSplit = []string{string(hashBytes)}
			if lineString != "" {
				lineSplit = append(lineSplit, strings.Split(lineString, " ")...)
			}
		}

		if dirName != "" {
			hashString := hex.EncodeToString([]byte(lineSplit[0]))
			hash, err := sha.ReadHash(hashString)
			if err != nil {
				return nil, err
			}
			treeObject, err := GetObject(rootGoitPath, hash)
			if err != nil {
				return nil, err
			}
			var path string
			if rootDir == "" {
				path = dirName
			} else {
				path = fmt.Sprintf("%s/%s", rootDir, dirName)
			}
			getEntries, err := treeObject.ExtractEntries(rootGoitPath, path)
			if err != nil {
				return nil, err
			}
			entries = append(entries, getEntries...)
			dirName = ""
		}

		if filePath != "" {
			hashString := hex.EncodeToString([]byte(lineSplit[0]))
			hash, err := sha.ReadHash(hashString)
			if err != nil {
				return nil, err
			}
			entry := index.NewEntry(hash, []byte(filePath))
			entries = append(entries, entry)
			filePath = ""
		}

		if len(lineSplit) == 1 { // last line
			break
		} else {
			var fileMode string
			var fileName string
			if len(lineSplit) == 2 {
				fileMode = lineSplit[0]
				fileName = lineSplit[1]
			} else if len(lineSplit) == 3 {
				fileMode = lineSplit[1]
				fileName = lineSplit[2]
			}

			if fileMode == "040000" {
				dirName = fileName
			} else if fileMode == "100644" {
				if rootDir == "" {
					filePath = fileName
				} else {
					filePath = fmt.Sprintf("%s/%s", rootDir, fileName)
				}
			}
		}
	}
	return entries, nil
}

func (t *tree) String() string {
	var lines []string

	for _, childNode := range t.children {
		var line string
		if len(childNode.children) == 0 {
			line = fmt.Sprintf("100644 blob %s\t%s", childNode.hash, childNode.name)
		} else {
			line = fmt.Sprintf("040000 tree %s\t%s", childNode.hash, childNode.name)
		}
		lines = append(lines, line)
	}

	message := strings.Join(lines, "\n")
	return message
}
