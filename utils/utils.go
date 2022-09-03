package utils

import (
	"strings"
	"unicode"
)

func GetPageIdFromUrl(pageId string) string {
	items := strings.Split(pageId, "notion.so/")
	if len(items) > 1 {
		return items[len(items)-1]
	}

	return pageId
}

// PurePageId find page id has "-" eg. cloud-native-f76464d369804e42bb67b180e7155a11
// need remove the prefix return f76464d369804e42bb67b180e7155a11
func PurePageId(pageId string) string {
	if strings.Index(pageId, "-") >= 0 {
		items := strings.Split(pageId, "-")
		pageId = items[len(items)-1]
	}

	return pageId
}

// PageUuid transfer page id eg. Istio-f76464d369804e42bb67b180e7155a11 to f76464d3-6980-4e42-bb67-b180e7155a11
func PageUuid(id string) (uuid string) {
	id = PurePageId(id)
	// make no side effects
	if strings.Index(id, "-") >= 0 {
		return id
	}

	return strings.Join([]string{id[:8], id[8:12], id[12:16], id[16:20], id[20:32]}, "-")
}

func UuidToPageId(uuid string) (pageId string) {
	items := strings.Split(uuid, "-")
	return strings.Join(items, "")
}

func PathPath(t string) string {
	t = strings.TrimSpace(t)
	r := []rune(t)
	nr := make([]rune, 0, 5)
	for _, v := range r {
		if unicode.IsDigit(v) || unicode.IsLetter(v) {
			nr = append(nr, unicode.ToLower(v))
		}

		if unicode.IsSpace(v) {
			nr = append(nr, '-')
		}
	}

	return string(nr)
}
