package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	fdir := `D:\shajf_dev\self\GBWBasBot\detect\src`

	filepath.Walk(fdir, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {

			fmt.Println(info.Name())
		}

		return nil
	})
}
