package meta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func LoadPageChunk(pageId string) (pageChunkResp *PageChunkResponse, err error) {
	// find page id has "-" eg. cloud-native-f76464d369804e42bb67b180e7155a11
	// need remove the prefix
	if strings.Index(pageId, "-") >= 0 {
		items := strings.Split(pageId, "-")
		pageId = items[len(items)-1]
	}

	req := &PageChunkRequest{
		PageID:          ToUuid(pageId),
		Limit:           50,
		ChunkNumber:     1,
		VerticalColumns: false,
	}

	reqStr, _ := json.Marshal(req)
	fmt.Println(string(reqStr))

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

	err = json.Unmarshal(body, &pageChunkResp)
	if err != nil {
		fmt.Println(string(body), err)
		return nil, err
	}

	return
}
