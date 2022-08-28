package utils

import "strings"

// PurePageId find page id has "-" eg. cloud-native-f76464d369804e42bb67b180e7155a11
// need remove the prefix
func PurePageId(pageId string) string {
	if strings.Index(pageId, "-") >= 0 {
		items := strings.Split(pageId, "-")
		pageId = items[len(items)-1]
	}

	return pageId
}
