package meta

type PageChunkRequest struct {
	PageID          string `json:"pageId"`
	Limit           int    `json:"limit"`
	ChunkNumber     int    `json:"chunkNumber"`
	VerticalColumns bool   `json:"verticalColumns"`
}
