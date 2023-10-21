## Hexo2Hugo

## 说明

使用 `golang` 原生编写,将 `hexo` 文章的 `front matter` 标签改为 `hugo` 可以解析的格式.

过滤了常用的 `date tags categories excerpt` 标签, 如有其他定义的标签,可以自己在源码添加

## 使用方法
`go run main.go`

## 支持的格式
目前支持两种格式
- 单行标签：如：
  ```
    tags: hexo,blog
  ```  
- 多行标签
```
    tags: 
        - hexo
        - blog
```

## 参考
1. [Front matter | Hugo](https://gohugo.io/content-management/front-matter/)
2. [ayuayue/hexo2hugo: hexo's front mater to hugo](https://github.com/ayuayue/hexo2hugo)
3. [Hugo中的Front Matter与Archetypes · 零壹軒·笔记](https://note.qidong.name/2017/06/21/hugo_front_matter/)