package meta

import (
	"bytes"
	"encoding/json"
	"github.com/foreversmart/notion_blog/log"
	"github.com/foreversmart/notion_blog/utils"
	"io/ioutil"
	"net/http"
)

func LoadPageChunk(pageId string) (pageChunkResp *PageChunkResponse, err error) {
	pageId = utils.PurePageId(pageId)

	req := &PageChunkRequest{
		PageID:      utils.PageUuid(pageId),
		Limit:       50,
		ChunkNumber: 0,
		Cursor: &PageCursor{
			Stack: [][]*PageStackItem{{
				{
					Table: "block",
					ID:    utils.PageUuid(pageId),
					Index: 0,
				},
			}},
		},
		VerticalColumns: false,
	}

	reqStr, _ := json.Marshal(req)
	//log.Logger.Info(string(reqStr))

	request, _ := http.NewRequest(http.MethodPost, "https://www.notion.so/api/v3/loadPageChunk", bytes.NewReader(reqStr))
	request.Header.Set("content-type", "application/json")
	request.Header.Set("accept-language", "en,zh-CN;q=0.9,zh;q=0.8")
	request.Header.Set("referer", "https://www.notion.so/"+pageId)
	request.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//log.Logger.Infof("meta response: %s", string(body))

	err = json.Unmarshal(body, &pageChunkResp)
	if err != nil {
		log.Logger.Error(string(body), err)
		return nil, err
	}

	return
}
