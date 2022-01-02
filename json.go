package springcfg

import "encoding/json"

func JsonPrettyFmt(i interface{}) string {
	jsonStr, _ := json.MarshalIndent(i, "", "    ")
	return string(jsonStr)
}
