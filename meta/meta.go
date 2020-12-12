package meta

import (
	"fmt"
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

	PageCover string `json:"page_cover"`
}

func PageMetaInfo(pageId string) (meta *PageMeta, err error) {
	meta = &PageMeta{}
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
			fmt.Println("block:", block.Value.ID, block.Value.ParentID, block.Value.CreatedTime, block.Value.Properties.Title[0][0])
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

			// add meta properties
			if strings.HasPrefix(comment.Value.Text[0][0], "meta:") {
				s := strings.TrimPrefix(comment.Value.Text[0][0], "meta:")
				items := strings.Split(s, ":")
				if len(items) > 1 {
					meta.Meta[items[0]] = items[1]
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
