package meta

import "strings"

func ToUuid(id string) (uuid string) {
	return strings.Join([]string{id[:8], id[8:12], id[12:16], id[16:20], id[20:32]}, "-")
}
