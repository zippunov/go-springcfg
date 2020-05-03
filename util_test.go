package springcfg

import (
	"testing"
)

func TestYamlDocsMerge1(t *testing.T) {
	m1 := map[string]interface{}{
		"key1": 3,
		"key2": "val1",
		"key3": map[string]interface{}{
			"key4": true,
		},
	}
	m2 := map[string]interface{}{
		"key5": 3,
		"key1": false,
		"key3": map[string]interface{}{
			"key5": "val2",
		},
	}
	m3 := DeepMerge(m1, m2)
	if m3["key1"] != false {
		t.Errorf("Expected key1 to be false, got %v \n\n", m3["key1"])
	}
	if m3["key2"] != "val1" {
		t.Errorf("Expected key2 to be \"val1\", got %v \n\n", m3["key2"])
	}
	if m3["key5"] != 3 {
		t.Errorf("Expected key5 to be 3, got %v \n\n", m3["key5"])
	}
	childmap := m3["key3"].(map[string]interface{})
	if childmap["key4"] != true {
		t.Errorf("Expected key3.key4 to be true, got %v \n\n", childmap["key4"])
	}
	if childmap["key5"] != "val2" {
		t.Errorf("Expected key3.key5 to be \"val2\", got %v \n\n", childmap["key5"])
	}
}
