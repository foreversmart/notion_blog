package main

import (
	"encoding/json"
	"github.com/foreversmart/notion_blog/blog"
	"github.com/foreversmart/notion_blog/log"
	"github.com/foreversmart/notion_blog/meta"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type BlogConfig struct {
	OutPutPath string           `json:"out_put_path"`
	PageIds    []string         `json:"page_ids"`
	PageUpdate map[string]int64 `json:"page_info"`
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configContent, err := ioutil.ReadFile(filepath.Join(pwd, "config.json"))
	if err != nil {
		panic(err)
	}

	var config *BlogConfig
	err = json.Unmarshal(configContent, &config)
	if err != nil {
		panic(err)
	}

	for index, pageId := range config.PageIds {
		log.Logger.Info(pageId, index, len(config.PageIds))
		pageId = getPageId(pageId)
		pageMeta, err := meta.PageMetaInfo(pageId)
		log.Logger.Info(pageMeta)
		if err != nil {
			log.Logger.Error("page", pageId, err)
			continue
		}

		//fmt.Println(pageMeta)

		if update, ok := config.PageUpdate[pageId]; ok {
			if pageMeta.LastModifyTimestamp < update {
				// 已经有更新记录且页面最近没有修改过的不更新
				continue
			}
		}

		// generate hugo blog to output target
		content, err := blog.NewBlog().HugoBlog(pageId, pageMeta)
		if err != nil {
			log.Logger.Error("HugoBlog", err)
		} else {
			ioutil.WriteFile(filepath.Join(config.OutPutPath, pageId+".html"), []byte(content), os.ModePerm)
		}

	}

}

func getPageId(pageId string) string {
	items := strings.Split(pageId, "notion.so/")
	if len(items) > 1 {
		return items[len(items)-1]
	}

	return pageId
}
