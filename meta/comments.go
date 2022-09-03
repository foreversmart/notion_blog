package meta

import "strings"

func CommentMetaAndTags(comments []string) (metaMap map[string]string, tags []string) {
	metaMap = make(map[string]string)
	for _, comment := range comments {
		// 新版本 notion 融合多个 comment 到一个 comment 为一行
		lines := strings.Split(comment, "\n")
		for _, line := range lines {
			// add meta properties
			if strings.HasPrefix(line, "meta:") {
				s := strings.TrimPrefix(line, "meta:")
				items := strings.Split(s, ":")
				if len(items) > 1 {
					metaMap[items[0]] = items[1]
				}
			}

			// add tags properties
			if strings.HasPrefix(line, "tags:") {
				if strings.HasPrefix(line, "tag:") {
					s := strings.TrimPrefix(line, "tag:")
					items := strings.Split(s, "/")
					tags = append(tags, items...)
				}

				if strings.HasPrefix(line, "tags:") {
					s := strings.TrimPrefix(line, "tags:")
					items := strings.Split(s, "/")
					tags = append(tags, items...)
				}
			}

		}
	}
	return
}
