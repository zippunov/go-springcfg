package springcfg

import "encoding/json"

func JSONPrettyPrint(i interface{}) string {
	jsonStr, _ := json.MarshalIndent(i, "", "    ")
	return string(jsonStr)
}
