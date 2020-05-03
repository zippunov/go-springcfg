package springcfg

import (
	"bufio"
	"io"
	"strings"

	"github.com/spf13/cast"
	yaml "gopkg.in/yaml.v2"
)

func unmarshalDoc(body string) (map[string]interface{}, error) {
	m := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(body), &m); err != nil {
		return nil, err
	}
	result := toStringMap(m)
	result = expandProps(result)
	return result.(map[string]interface{}), nil
}

func toStringMap(i interface{}) interface{} {
	switch v := i.(type) {
	case map[interface{}]interface{}:
		var m = map[string]interface{}{}
		for k, val := range v {
			m[cast.ToString(k)] = toStringMap(val)
		}
		return m
	default:
		return i
	}
}

func expandProps(i interface{}) interface{} {
	switch m := i.(type) {
	case map[string]interface{}:
		newM := map[string]interface{}{}
		for k, val := range m {
			path := buildPath(k, ".")
			if len(path) > 1 {
				lowK := path[0]
				highK := strings.Join(path[1:], ".")
				if _, ok := newM[lowK]; ok {
					newM[lowK].(map[string]interface{})[highK] = val
				} else {
					newM[lowK] = map[string]interface{}{
						highK: val,
					}
				}
			} else {
				newM[k] = val
			}
		}
		for k, val := range newM {
			newM[k] = expandProps(val)
		}
		return newM
	default:
		return i
	}
}

func fetchDocs(reader io.Reader) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	scanner := bufio.NewScanner(reader)
	doclines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "---") {
			body := strings.Join(doclines, "\n")
			doc, err := unmarshalDoc(body)
			if err != nil {
				return nil, err
			}
			doclines = doclines[:0]
			result = append(result, doc)
		} else {
			doclines = append(doclines, line)
		}
	}
	if len(doclines) > 0 {
		body := strings.Join(doclines, "\n")
		doc, err := unmarshalDoc(body)
		if err != nil {
			return nil, err
		}
		result = append(result, doc)
	}
	return result, nil
}
