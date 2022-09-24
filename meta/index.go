package meta

import (
	"fmt"
	"github.com/foreversmart/notion_blog/utils"
)

func WalkIndex(entrances []string) (pages map[string]*PageMeta, err error) {
	pages = make(map[string]*PageMeta)

	// record page has already visited
	pageMap := make(map[string]bool)
	// page chan queue
	pageChan := make(chan string, 10000)
	// page folders map
	pageFolderMap := make(map[string][]string)

	// in queue
	for _, e := range entrances {
		e = utils.PurePageId(e)
		pageChan <- e
		pageMap[e] = true
	}

	for len(pageChan) > 0 {
		pageId := <-pageChan
		fmt.Println("processing", pageId)
		pageChunk, err := LoadPageChunk(pageId)
		if err != nil {
			return nil, err
		}

		pageFolders := pageFolderMap[pageId]
		pageMeta := pageChunk.CurrentPageMeta(pageId)
		fmt.Println(".......", pageMeta)

		// deal if page is index page
		if pageMeta.Meta["index"] == "index" {
			subPages := pageChunk.SubPages(pageId)
			pageFolders = append(pageFolders, pageMeta.Title)

			for sId, s := range subPages {
				sId = utils.UuidToPageId(sId)
				if len(s.Meta["sub_title"]) > 0 {
					s.PageFolders = pageFolders
					pages[sId] = s
					continue
				}

				if s.Meta["index"] == "index" && !pageMap[sId] {
					pageChan <- sId
					pageMap[sId] = true
					pageFolderMap[sId] = pageFolders
				}
			}

			continue
		}

		// normal page
		if len(pageMeta.Meta["sub_title"]) > 0 {
			pageMeta.PageFolders = pageFolders
			pages[pageId] = pageMeta
			continue
		}

	}

	return
}
