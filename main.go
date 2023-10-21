package main

import (
	"hexo2hugo/file"
)

func main() {
	dir := "./source" //当前目录
	file.GetAllMDFileNames(dir)
	// file.ReadAllMDFile(filePathNames)
}
