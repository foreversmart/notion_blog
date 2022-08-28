package blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/foreversmart/notion_blog/lib/bdfanyi"
	"github.com/foreversmart/notion_blog/log"
	"github.com/foreversmart/notion_blog/meta"
	"github.com/go-rod/rod"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type HugoBlogMeta struct {
	Title       string   `json:"title"`
	SubTitle    string   `json:"sub_title"`
	Excerpt     string   `json:"excerpt"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Image       string   `json:"image"`
	Tags        []string `json:"tags"`
	Category    []string `json:"category"`
	UrlPath     string   `json:"url_path"`
	Url         string   `json:"url"`
}

type Blog struct {
	Browser *rod.Browser
}

func NewBlog() *Blog {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New()
	err := browser.Connect()
	if err != nil {
		panic(err)
	}

	browser = browser.Timeout(time.Minute)

	return &Blog{
		Browser: browser,
	}
}

func (b *Blog) Close() error {
	return b.Browser.Close()
}

func (b *Blog) PageIndex(pageId string) (subPages []string, err error) {
	url := Host + pageId

	page := b.Browser.MustPage(url)
	// sleep enough time to wait page render
	time.Sleep(time.Second * 7)

	eles, err := page.Elements(".notion-page-block")
	if err != nil {
		return nil, err
	}

	for _, ele := range eles {
		aEle, err := ele.Element("a")
		if err != nil {
			continue
		}

		aHref, err := aEle.Attribute("href")
		if err != nil {
			continue
		}

		subPages = append(subPages, *aHref)
	}

	page.Close()

	return
}

func (b *Blog) WalkPages(entrances []string) (allPages []string, err error) {
	// record page has already visited
	pageMap := make(map[string]bool)
	// page chan queue
	pageChan := make(chan string, 10000)

	// in queue
	for _, e := range entrances {
		pageChan <- e
		pageMap[e] = true
	}

	for len(pageChan) > 0 {
		pageId := <-pageChan
		pageMeta, err := meta.PageMetaInfo(pageId)
		if err != nil {
			return nil, err
		}

		// if is index page
		if pageMeta.Meta["index"] == "index" {
			subPages, err := b.PageIndex(pageId)
			if err != nil {
				return nil, err
			}

			for _, s := range subPages {
				if !pageMap[s] {
					pageChan <- s
					pageMap[s] = true
				}
			}

		}

		// find the publish page
		if len(pageMeta.Meta["sub_title"]) > 0 {
			allPages = append(allPages, pageId)
		}
	}

	return
}

func (b *Blog) HugoBlog(pageId string, pageMeta *meta.PageMeta) (content string, err error) {
	url := Host + pageId

	page := b.Browser.MustPage(url)

	// sleep enough time to wait page render
	time.Sleep(time.Second * 7)

	title, _, err := PageTile(page)
	if err != nil {
		return "", err
	}
	pageContent, err := PageContent(pageId, page)
	if err != nil {
		return "", err
	}

	category := make([]string, 0, 5)
	if len(pageMeta.Titles) > 0 {
		//pageMeta.Tags = append(pageMeta.Tags, pageMeta.Titles[1])
		category = append(category, pageMeta.Titles[1])
	}
	category = append(category, PageCategory(pageMeta)...)

	subTitle := HugoPageUrl(title, pageId)
	desc, err := PageDesc(page)
	if err != nil {
		log.Logger.Error(err)
		return "", err
	}

	r := []rune(desc)
	if len(r) > 150 {
		desc = strconv.Quote(string(r[0:150]))
		desc = strings.Trim(desc, `""`)
	}

	hugoMeta := &HugoBlogMeta{
		Title:       pageMeta.Title,
		SubTitle:    subTitle,
		Excerpt:     subTitle,
		Description: desc,
		Date:        pageMeta.CreatedAt,
		Image:       pageMeta.PageCover,
		Tags:        pageMeta.Tags,
		Category:    category,
		UrlPath:     category[0],
		Url:         subTitle,
	}

	ts, _ := json.Marshal(hugoMeta)
	log.Logger.Infof("page hugo meta %s", string(ts))
	// rewrite some meta from pageMeta
	if v, ok := pageMeta.Meta["title"]; ok {
		hugoMeta.Title = v
	}
	if v, ok := pageMeta.Meta["created_at"]; ok {
		hugoMeta.Date = v
	}
	if v, ok := pageMeta.Meta["sub_title"]; ok {
		hugoMeta.SubTitle = v
		// reset to new url
		hugoMeta.Url = v
	}
	if v, ok := pageMeta.Meta["desc"]; ok {
		hugoMeta.Description = v
	}
	if v, ok := pageMeta.Meta["url_path"]; ok {
		hugoMeta.UrlPath = v
	}
	if v, ok := pageMeta.Meta["url"]; ok {
		hugoMeta.Url = v
	}

	tpl, err := template.New("hugo").Parse(HugoBlogHeaderTpl)
	if err != nil {
		return title, err
	}

	writer := &bytes.Buffer{}
	err = tpl.Execute(writer, hugoMeta)
	if err != nil {
		log.Logger.Error(err)
		return title, err
	}

	writer.WriteString(Head)
	//writer.WriteString(titleHtml)
	writer.WriteString(Mid)
	writer.WriteString(pageContent)
	writer.WriteString(Tail)
	return writer.String(), err

}

func PageCategory(m *meta.PageMeta) (categories []string) {
	for _, comment := range m.Comment {
		if strings.HasPrefix(comment, "cate:") {
			comment = strings.TrimPrefix(comment, "cate:")
			items := strings.Split(comment, "/")
			categories = append(categories, items...)
		}
	}

	return
}

func PageTile(page *rod.Page) (title, titleHtml string, err error) {
	titleEles, err := page.Elements(".notion-page-block")
	if err != nil {
		return "", "", fmt.Errorf("page title ele not found %w", err)
	}

	for i, ele := range titleEles {
		switch i {
		case 0:

			//category, err = ele.Text()
			//if err != nil {
			//	return "", "",  fmt.Errorf("page category text not found %w", err)
			//}
			//
			//// filter category
			//items := strings.Split(category, "\n")
			//if len(items) > 1 {
			//	category = items[1]
			//}
		case 1:
			title, err = ele.Text()
			if err != nil {
				return "", "", fmt.Errorf("page title text not found %w", err)
			}

			titleHtml, err = ele.HTML()
			if err != nil {
				return "", "", fmt.Errorf("page title text not found %w", err)
			}

			// append title style
			titleHtml = TitleStyle + titleHtml + "</div>"

		}
	}

	return
}

func PageContent(pageId string, page *rod.Page) (content string, err error) {
	ele, err := page.Element(".notion-page-content")
	if err != nil {
		return "", nil
	}

	content = ele.MustHTML()
	// adjust .notion-page-content style is first style
	content = strings.Replace(content, "padding-bottom: 30vh;", "padding-bottom: 10vh;", 1)
	content = strings.Replace(content, "padding-left: calc(96px + env(safe-area-inset-left));", "padding-left: 5%;", 1)
	content = strings.Replace(content, "padding-right: calc(96px + env(safe-area-inset-right));", "padding-right: 5%; margin: auto;", 1)
	content = strings.Replace(content, "padding-bottom: 30vh;", "padding-bottom: 10vh;", 1)
	// content filter
	content = strings.Replace(content, `contenteditable="true"`, "", -1)
	// convert local image to notion cdn image
	content = strings.Replace(content, `img src="/image`, `img src="https://www.notion.so/image`, -1)

	// replace html anchor
	content = strings.Replace(content, `data-block-id`, "id", -1)

	// match href="/20ecc4449a3843e6bea9324c3328241e#026465ad4766428784aa1163c89432f1" and replace to href="#026465ad-47664-28784aa11-63c89432f1"
	regStr := "href=\"/" + pageId + `#[\w]+"`
	reg := regexp.MustCompile(regStr)
	content = reg.ReplaceAllStringFunc(content, func(s string) string {
		if items := strings.Split(s, "#"); len(items) == 2 {
			return `href="` + "#" + meta.ToUuid(items[1]) + `"`
		}
		return s
	})

	return
}

func PageDesc(page *rod.Page) (content string, err error) {
	ele, err := page.Element(".notion-page-content")
	if err != nil {
		return "", nil
	}

	content, err = ele.Text()
	return
}

func HugoPageUrl(title string, pageId string) string {
	var tk, err = bdfanyi.Gtk()
	if err != nil {
		log.Logger.Error(err)
		return pageId
	}

	options := bdfanyi.NewOptions(bdfanyi.ZH, bdfanyi.EN, tk, "")
	r, err := bdfanyi.Do(title, options)
	if err != nil {
		log.Logger.Error(err)
		return pageId
	}

	return r.TransResult.Data[0].Dst
}
