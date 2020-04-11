package util

import "strings"

func Comment2Map(comment string) map[string]string {
	res := make(map[string]string)

	kvList := strings.Split(comment, "|")
	for _, kv := range kvList {
		v := strings.Split(kv, ":")
		if len(v) < 2 {
			continue
		}
		res[v[0]] = v[1]
	}

	return res
}
