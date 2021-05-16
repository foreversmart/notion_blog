package meta

import (
	"github.com/foreversmart/notion_blog/log"
	"math"
	"strings"
	"time"
)

type PageMeta struct {
	CreatedAt           string            `json:"created_at"`
	CreateTimestamp     int64             `json:"create_timestamp"`
	LastModifyAt        string            `json:"last_modify_at"`
	LastModifyTimestamp int64             `json:"last_modify_timestamp"`
	Titles              []string          `json:"titles"` // 0 is child and current and then is parent
	Title               string            `json:"title"`
	Comment             []string          `json:"comment"`
	Meta                map[string]string `json:"meta"`
	Tags                []string          `json:"tags"`

	PageCover string `json:"page_cover"`
}

func PageMetaInfo(pageId string) (meta *PageMeta, err error) {
	meta = &PageMeta{}
	meta.Meta = make(map[string]string)
	pageChunkResp, err := LoadPageChunk(pageId)
	if err != nil {
		return nil, err
	}
	// origin CreatedTime
	var createTime int64 = math.MaxInt64
	var lastModifyTime int64 = math.MaxInt64

	pageMapBlock := make(map[string]*PageChunkBlock)
	for _, block := range pageChunkResp.RecordMap.Block {
		if block.Value.Type == "page" {
			log.Logger.Debug("block:", block.Value.ID, block.Value.ParentID, block.Value.CreatedTime, block.Value.Properties.Title[0][0])
			pageMapBlock[block.Value.ID] = block
		}
		if createTime > block.Value.CreatedTime {
			createTime = block.Value.CreatedTime
		}

		if lastModifyTime > block.Value.LastEditedTime {
			lastModifyTime = block.Value.LastEditedTime
		}
	}

	// deal comment
	for _, comment := range pageChunkResp.RecordMap.Comment {
		if len(comment.Value.Text) > 0 && len(comment.Value.Text[0]) > 0 {
			meta.Comment = append(meta.Comment, comment.Value.Text[0][0])
			// 新版本 notion 融合多个 comment 到一个 comment 为一行
			lines := strings.Split(comment.Value.Text[0][0], "\n")
			for _, line := range lines {
				// add meta properties
				if strings.HasPrefix(line, "meta:") {
					s := strings.TrimPrefix(line, "meta:")
					items := strings.Split(s, ":")
					if len(items) > 1 {
						meta.Meta[items[0]] = items[1]
					}
				}

				// add tags properties
				if strings.HasPrefix(line, "tags:") {
					if strings.HasPrefix(line, "tag:") {
						s := strings.TrimPrefix(line, "tag:")
						items := strings.Split(s, "/")
						meta.Tags = append(meta.Tags, items...)
					}
					if strings.HasPrefix(line, "tags:") {
						s := strings.TrimPrefix(line, "tags:")
						items := strings.Split(s, "/")
						meta.Tags = append(meta.Tags, items...)
					}
				}

			}

		}

	}

	hasChildPageMap := make(map[string]bool) // record page map block which has no child page
	// find current page block
	var currentBlock *PageChunkBlock
	for _, block := range pageMapBlock {
		if _, ok := pageMapBlock[block.Value.ParentID]; ok {
			hasChildPageMap[block.Value.ParentID] = true
		}
	}

	for id, _ := range pageMapBlock {
		if ok := hasChildPageMap[id]; !ok {
			currentBlock = pageMapBlock[id]
		}
	}

	if currentBlock != nil {
		meta.PageCover = currentBlock.SourceImageUrl(currentBlock.Value.Format.PageCover)
		title, _ := currentBlock.Value.Properties.Title[0][0].(string)
		meta.Title = title
	}

	for {
		if currentBlock == nil {
			break
		}

		title, _ := currentBlock.Value.Properties.Title[0][0].(string)
		meta.Titles = append(meta.Titles, title)
		currentBlock = pageMapBlock[currentBlock.Value.ParentID]
		if currentBlock == nil {
			break
		}
	}

	//originBlock, ok := pageChunkResp.RecordMap.Block[ToUuid(pageId)]
	//if ok {
	//	for _, k := range originBlock.Value.
	//}

	createTimeDate := time.Unix(createTime/1000, 0)
	meta.CreateTimestamp = createTime
	meta.CreatedAt = createTimeDate.Format("2006-01-02")
	lastModifyTimeDate := time.Unix(lastModifyTime/1000, 0)
	meta.LastModifyTimestamp = lastModifyTime
	meta.LastModifyAt = lastModifyTimeDate.Format("2006-01-02")
	return
}
