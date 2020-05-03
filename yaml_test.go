package springcfg

import (
	"strings"
	"testing"
)

func TestFetchMultipleYAMLDocs(t *testing.T) {
	var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
  e:
    a: tra-ta-ta
    b: 345
    c: True
---
e:
  - 1
  - 2
`
	reader := strings.NewReader(data)
	docs, err := fetchDocs(reader)
	if err != nil {
		t.Error(err)
	}
	if len(docs) != 2 {
		t.Errorf("Got %v docs, expected 2", len(docs))
	}
}
