package meta

import (
	"net/url"
	"strings"
)

type PageChunkBlock struct {
	Role  string `json:"role"`
	Value struct {
		ID         string     `json:"id"`
		Version    int        `json:"version"`
		Type       string     `json:"type"`
		Text       [][]string `json:"text"` // comment
		Properties struct {
			Title [][]interface{} `json:"title"` // most string otherwise [[]] style params
		} `json:"properties"`
		Format struct {
			PageIcon          string  `json:"page_icon"`
			PageCover         string  `json:"page_cover"`
			PageFullWidth     bool    `json:"page_full_width"`
			PageCoverPosition float64 `json:"page_cover_position"`
		}
		CreatedTime       int64  `json:"created_time"`
		LastEditedTime    int64  `json:"last_edited_time"`
		ParentID          string `json:"parent_id"`
		ParentTable       string `json:"parent_table"`
		Alive             bool   `json:"alive"`
		CreatedByTable    string `json:"created_by_table"`
		CreatedByID       string `json:"created_by_id"`
		LastEditedByTable string `json:"last_edited_by_table"`
		LastEditedByID    string `json:"last_edited_by_id"`
		ShardID           int    `json:"shard_id"`
		SpaceID           string `json:"space_id"`
	} `json:"value"`
}

func (block *PageChunkBlock) SourceImageUrl(s string) string {
	base := "https://www.notion.so"
	originUrl := ""
	if strings.HasPrefix(s, "http") {
		// cdn pics
		originUrl = s
	} else {
		// notion internal images
		originUrl = base + s
	}
	originUrl = url.PathEscape(originUrl)
	return base + "/image/" + originUrl + "?table=block&id=" + block.Value.ID + "&width=3840&userId=&cache=v2"
}
