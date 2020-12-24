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
	OutPutPath string                 `json:"out_put_path"`
	PageIds    []string               `json:"page_ids"`
	PageConfig map[string]*PageConfig `json:"page_config"`
	PageUpdate map[string]int64       `json:"page_info"`
}

type PageConfig struct {
	*meta.PageMeta
	IsRender bool `json:"is_render"`
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(pwd, "config.json")
	configContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	var config *BlogConfig
	err = json.Unmarshal(configContent, &config)
	if err != nil {
		panic(err)
	}

	for index, pageId := range config.PageIds {
		log.Logger.Info(pageId, "index:", index, len(config.PageIds))

		pageId = getPageId(pageId)
		if v, ok := config.PageConfig[pageId]; ok {
			if v.IsRender == false {
				continue
			}
		}

		pageMeta, err := meta.PageMetaInfo(pageId)
		log.Logger.Info(pageMeta)
		if err != nil {
			log.Logger.Error("page", pageId, err)
			continue
		}

		config.PageConfig[pageId] = &PageConfig{
			PageMeta: pageMeta,
			IsRender: true,
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

	configNew, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Logger.Errorf("%v", err)
		return
	}

	err = ioutil.WriteFile(configFile, configNew, os.ModePerm)
	if err != nil {
		log.Logger.Error(err)
	}

}

func getPageId(pageId string) string {
	items := strings.Split(pageId, "notion.so/")
	if len(items) > 1 {
		return items[len(items)-1]
	}

	return pageId
}
