package main

import (
	"encoding/json"
	"github.com/foreversmart/notion_blog/blog"
	"github.com/foreversmart/notion_blog/log"
	"github.com/foreversmart/notion_blog/meta"
	"github.com/foreversmart/notion_blog/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type BlogConfigV1 struct {
	OutPutPath string   `json:"out_put_path"`
	Entrances  []string `json:"entrances"`
}

//
//type PageCache struct {
//}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(pwd, "config.local.json")
	configContent, err := ioutil.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	} else if os.IsNotExist(err) {
		// config.local.json not exist fail back to config.json
		configFile = filepath.Join(pwd, "config.json")
		configContent, err = ioutil.ReadFile(configFile)
		if err != nil {
			panic(err)
		}
	}

	var config *BlogConfigV1
	err = json.Unmarshal(configContent, &config)
	if err != nil {
		panic(err)
	}

	for i, v := range config.Entrances {
		config.Entrances[i] = utils.GetPageIdFromUrl(v)
	}

	//subPages, err := meta.WalkIndex([]string{"blog-44c506d8d2b84dfd818236cd14075410"})
	subPages, err := meta.WalkIndex(config.Entrances)
	if err != nil {
		panic(err)
	}

	for pageId, pageMeta := range subPages {
		//fmt.Println(pageId)
		//fmt.Println(pageMeta)
		pageFileName := filepath.Join(config.OutPutPath, pageMeta.PageName+".html")
		info, err := os.Stat(pageFileName)
		if err != nil {
			if os.IsNotExist(err) {
				// create
				log.Logger.Info("create new page:", pageFileName)
				err = CreatePageFile(pageFileName, pageId, pageMeta)
				if err != nil {
					log.Logger.Error(err)
				}
				continue
			}
			panic(err)
		}

		if info.ModTime().Before(time.Unix(pageMeta.LastModifyTimestamp/1000, 0)) {
			log.Logger.Infof("update page: %s modify time %v and %v", pageFileName, info.ModTime(), time.Unix(pageMeta.LastModifyTimestamp, 0))
			err = CreatePageFile(pageFileName, pageId, pageMeta)
			if err != nil {
				log.Logger.Error(err)
			}
		}

	}
}

func CreatePageFile(filename, pageId string, pageMeta *meta.PageMeta) (err error) {
	log.Logger.Info("create page file", filename, pageId)
	// create
	content, err := blog.NewBlog().HugoBlog(pageId, pageMeta)
	if err != nil {
		log.Logger.Error("HugoBlog", err)
		return
	}
	return ioutil.WriteFile(filename, []byte(content), os.ModePerm)
}
