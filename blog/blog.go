package blog

import (
	"bytes"
	"fmt"
	"github.com/foreversmart/notion_blog/log"
	"github.com/foreversmart/notion_blog/meta"
	"github.com/go-rod/rod"
	bdfanyi "github.com/hnmaonanbei/go-baidu-fanyi"
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

func (b *Blog) HugoBlog(pageId string, pageMeta *meta.PageMeta) (content string, err error) {
	url := Host + pageId
	page := b.Browser.MustPage(url)

	// sleep enough time to wait page render
	time.Sleep(time.Second * 5)

	title, _, err := PageTile(page)
	if err != nil {
		return "", err
	}
	pageContent, err := PageContent(page)
	if err != nil {
		return "", err
	}

	tags := make([]string, 0, 5)
	category := make([]string, 0, 5)
	if len(pageMeta.Titles) > 0 {
		tags = append(tags, pageMeta.Titles[1])
		category = append(category, pageMeta.Titles[1])
	}
	tags = append(tags, PageTags(pageMeta)...)
	category = append(category, PageCategory(pageMeta)...)

	subTitle := HugoPageUrl(title, pageId)
	desc, err := PageDesc(page)
	if err != nil {
		log.Logger.Error(err)
		return "", err
	}

	r := []rune(desc)
	if len(r) > 150 {
		desc = string(r[0:150])
	}

	hugoMeta := &HugoBlogMeta{
		Title:       title,
		SubTitle:    subTitle,
		Excerpt:     subTitle,
		Description: desc,
		Date:        pageMeta.CreatedAt,
		Image:       pageMeta.PageCover,
		Tags:        tags,
		Category:    category,
		UrlPath:     category[0],
		Url:         subTitle,
	}

	// rewrite some meta from pageMeta
	if v, ok := pageMeta.Meta["title"]; ok {
		hugoMeta.Title = v
	}
	if v, ok := pageMeta.Meta["created_at"]; ok {
		hugoMeta.Date = v
	}
	if v, ok := pageMeta.Meta["sub_title"]; ok {
		hugoMeta.SubTitle = v
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

func PageTags(m *meta.PageMeta) (tags []string) {
	for _, comment := range m.Comment {
		if strings.HasPrefix(comment, "tag:") {
			comment = strings.TrimPrefix(comment, "tag:")
			items := strings.Split(comment, "/")
			tags = append(tags, items...)
		}
		if strings.HasPrefix(comment, "tags:") {
			comment = strings.TrimPrefix(comment, "tags:")
			items := strings.Split(comment, "/")
			tags = append(tags, items...)
		}
	}

	return
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

func PageComment(page *rod.Page) () {

}

func PageContent(page *rod.Page) (content string, err error) {
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
	tk, err := bdfanyi.Gtk()
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
