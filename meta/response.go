package meta

import (
	"github.com/foreversmart/notion_blog/log"
	"github.com/foreversmart/notion_blog/utils"
	"time"
)

type PageChunkResponse struct {
	RecordMap struct {
		Block      map[string]*PageChunkBlock `json:"block"`
		Comment    map[string]*PageChunkBlock `json:"comment"`
		Discussion map[string]*PageDiscussion `json:"discussion"`
	} `json:"recordMap"`
}

// CurrentPageMeta fetch current page meta from response
// notion page may contains child page comments in parent page chunk
func (r *PageChunkResponse) CurrentPageMeta(pageId string) (meta *PageMeta) {

	pageUUID := utils.PageUuid(pageId)
	var currentBlock *PageChunkBlock

	pageMapBlock := make(map[string]*PageChunkBlock)
	for _, block := range r.RecordMap.Block {
		if block.Value.Type == "page" {
			log.Logger.Debug("block:", block.Value.ID, block.Value.ParentID, block.Value.CreatedTime, block.Value.Properties.Title[0][0])
			pageMapBlock[block.Value.ID] = block

			// current page block
			if block.Value.ID == pageUUID {
				currentBlock = block
			}
		}
	}

	meta = r.FetchBlockMeta(currentBlock)

	return
}

func (r *PageChunkResponse) FetchBlockMeta(block *PageChunkBlock) (meta *PageMeta) {
	meta = &PageMeta{}
	meta.Meta = make(map[string]string)
	if block == nil {
		return nil
	}

	// fetch page cover and title
	meta.PageCover = block.SourceImageUrl(block.Value.Format.PageCover)
	title, _ := block.Value.Properties.Title[0][0].(string)
	meta.Title = title

	// deal page block timestamp
	meta.CreateTimestamp = block.Value.CreatedTime
	meta.LastModifyTimestamp = block.Value.LastEditedTime
	createTimeDate := time.Unix(meta.CreateTimestamp/1000, 0)
	meta.CreatedAt = createTimeDate.Format("2006-01-02")
	lastModifyTimeDate := time.Unix(meta.LastModifyTimestamp/1000, 0)
	meta.LastModifyAt = lastModifyTimeDate.Format("2006-01-02")

	// add find current block comments and parse it to tag and meta
	for _, disId := range block.Value.Discussions {
		dis, ok := r.RecordMap.Discussion[disId]
		if !ok {
			continue
		}

		for _, commentId := range dis.Value.Comments {
			comment, ok := r.RecordMap.Comment[commentId]
			if !ok {
				continue
			}
			meta.Comment = append(meta.Comment, comment.RawComments()...)
		}

	}

	meta.Meta, meta.Tags = CommentMetaAndTags(meta.Comment)
	meta.SubTitle = meta.Meta["sub_title"]
	meta.PageName = utils.PathPath(meta.SubTitle)
	return
}

func (r *PageChunkResponse) SubPages(pageId string) (res map[string]*PageMeta) {
	res = make(map[string]*PageMeta)
	pageUUID := utils.PageUuid(pageId)

	for _, block := range r.RecordMap.Block {
		if block.Value.Type == "page" {
			// current page block
			if block.Value.ID == pageUUID {
				continue
			}

			childPageId := block.Value.ID
			res[childPageId] = r.FetchBlockMeta(block)
		}
	}

	return
}
