package file

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//GetAllMDFileName 获取当前目录下的所有指定格式的文件
func GetAllMDFileNames(dir string) []string {
	var filePathNames []string
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			filePathNames = append(filePathNames, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	count := len(filePathNames)
	if count < 1 {
		fmt.Println("指定目录下没有markdown文件")
	} else {
		fmt.Printf("指定目录及其子目录下共检测到 %d 个 markdown 文件\n", count)
	}

	// for i := range filePathNames {
	// 	fmt.Println(filePathNames[i])
	// }

	ReadAllMDFiles(filePathNames)

	return filePathNames
}


//ReadAllMDFile 读取所有获取到的文件
func ReadAllMDFiles(filePathNames []string) {
    newDir := fmt.Sprintf("modified%d", time.Now().Unix())
    err := os.Mkdir(newDir, 0755)
    if err != nil {
        log.Fatal(err)
    }

    // 遍历读取所有的文件
    for k, v := range filePathNames {
        file, err := os.Open(v)
        if err != nil {
            log.Printf("can't open file %s, err : %s ", v, err)
            continue // 继续处理下一个文件
        }
        fmt.Printf("\n正在读取第 %d 个文件 %s\n", k+1, v)
        result := ""
        scanner := bufio.NewScanner(file)

        inMultiLineTag := false
        multiLineTag := ""

        for scanner.Scan() {
            line := scanner.Text()
            // 如果行中包含 date 标签,则替换掉
            if strings.Contains(line, "date:") {
                line = HandleDate(line)
            }
			// 如果行中包含 “tags:” 标签,则开始检测是否有多行标签
            if strings.Contains(line, "tags:") {
                inMultiLineTag = true
                multiLineTag = line
                continue
            }
            if inMultiLineTag {
                if strings.HasPrefix(line, "  - ") {
                    multiLineTag += "\n" + line
					continue
                } else {
                    inMultiLineTag = false
					//fmt.Println(multiLineTag)
                    result = result + HandleTags(multiLineTag) + "\n"
                }
            }
            result = result + line + "\n"
        }
        fmt.Printf("正在更改第 %d 个文件 %s\n", k+1, v)

        if err := HandleContent(newDir, v, result); err != nil {
            log.Fatal(err)
            log.Fatal("写入文件失败")
        }
        fmt.Printf("%s写入完成\n", v)
    }
}


//HandleContent 将处理完的内容覆盖进去
func HandleContent(newDir, filePathNames, result string) error {
	if strings.Contains(filePathNames, "\\") {
		filePathNames = strings.Replace(filePathNames, "\\", "/", -1)
	}
	
	// 创建目录（包括中间目录）如果不存在
	err := os.MkdirAll(filepath.Dir(filepath.Join(newDir, filePathNames)), 0755)
	if err != nil {
		return err
	}

	// 创建或打开文件
	f, err := os.OpenFile(filepath.Join(newDir, filePathNames), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.WriteString(f, result)
	if err != nil {
		return err
	}
	return nil
}


//HandleTags 处理标签
// 单行标签的情况
// func HandleTags(tags string) string {
// 	//查找tags后的空格,将所有的 tag 用 [] 包裹起来
// 	index := strings.Index(tags, " ")
// 	tags = fmt.Sprintf("%s [%s]", tags[:index], tags[index+1:])
// 	return tags
// }
// 多行标签
func HandleTags(tags string) string {
	lines := strings.Split(tags, "\n")
	var extractedTags []string
	inTagsSection := false
	tagsFound := false

	for _, line := range lines {
		if strings.HasPrefix(line, "tags:") {
			inTagsSection = true
			tagsFound = true
		} else if inTagsSection && strings.HasPrefix(line, "  - ") {
			// Extract the tag (e.g., "  - 感想" => "感想")
			tag := strings.TrimSpace(strings.TrimPrefix(line, "  - "))
			extractedTags = append(extractedTags, tag)
		} else if inTagsSection && len(strings.TrimSpace(line)) == 0 {
			// End of the tags section
			break
		}
	}

	if tagsFound {
		// Format the tags as desired
		tagsStr := "tags = ["
		for i, tag := range extractedTags {
			if i > 0 {
				tagsStr += ", "
			}
			tagsStr += "\"" + tag + "\""
		}
		tagsStr += "]"
		return tagsStr
	}

	return tags
}

//HandleCategories 处理分类
func HandleCategories(categories string) string {
	//查找到分类的空格,整理成 hugo 的格式
	index := strings.Index(categories, " ")
	categories = fmt.Sprintf("%s [%s]", categories[:index], categories[index+1:])
	return categories
}

//HandleDate 处理日期
func HandleDate(date string) string {
	//查找到年月日后的空格,根据 hugo date 的格式进行更改
	index := strings.LastIndex(date, " ")
	date = fmt.Sprintf("%sT%s+08:00", date[:index], date[index+1:])
	return date
}
