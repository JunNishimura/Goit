package util

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/JunNishimura/Goit/object"
)

func GetFuncInput(filePath string, objType object.ObjectType) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("fail to read file: %v", err)
	}
	fileSize := len(bytes)
	input := string(objType) + " " + strconv.Itoa(fileSize) + "\x00" + string(bytes)
	return input, nil
}
