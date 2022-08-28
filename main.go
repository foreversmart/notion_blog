package main

import (
	"fmt"
	"github.com/foreversmart/notion_blog/blog"
)

func main() {
	subPages, err := blog.NewBlog().PageIndex("blog-44c506d8d2b84dfd818236cd14075410")
	if err != nil {
		panic(err)
	}

	fmt.Println(subPages)
}
