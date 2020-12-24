package meta

type PageChunkRequest struct {
	PageID          string      `json:"pageId"`
	Limit           int         `json:"limit"`
	ChunkNumber     int         `json:"chunkNumber"`
	Cursor          *PageCursor `json:"cursor"`
	VerticalColumns bool        `json:"verticalColumns"`
}

type PageCursor struct {
	Stack [][]*PageStackItem `json:"stack"`
}

type PageStackItem struct {
	Table string `json:"table"`
	ID    string `json:"id"`
	Index int    `json:"index"`
}
