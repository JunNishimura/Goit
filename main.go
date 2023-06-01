/*
Copyright © 2023 Jun Nishimura <n.junjun0303@gmail.com>
*/
package main

import (
	"fmt"

	"github.com/JunNishimura/Goit/cmd"
)

var version = ""

func main() {
	fmt.Println(version)
	cmd.Execute(version)
}
