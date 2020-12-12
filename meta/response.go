package meta

type PageChunkResponse struct {
	RecordMap struct {
		Block   map[string]*PageChunkBlock `json:"block"`
		Comment map[string]*PageChunkBlock `json:"comment"`
	} `json:"recordMap"`
}
