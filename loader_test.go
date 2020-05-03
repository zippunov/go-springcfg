package springcfg

import (
	"strings"
	"testing"
)

var loaderTestData = `
spring.application.name: api-chat
spring.cloud.config:
  enabled: true
  label: development
  uri: ${SPRING_CONFIG_URI:http://config-server:8888}
  vault_token: ${CONFIG_SERVER_VAULT_TOKEN}
---
profiles: development,travis
spring.cloud.config.enabled: false
---
profiles: qa,prod,staging,hotfix,demo
spring.cloud.config:
  enabled: true
  label: development
  uri: ${SPRING_CONFIG_URI:http://infra-config-server.service.consul:8888}
  vault_token: ${CONFIG_SERVER_VAULT_TOKEN}
`

func TestConfigLoaderShouldUseDoc(t *testing.T) {
	loader := Loader{
		ActiveProfiles: []string{
			"development",
			"local",
		},
	}
	reader := strings.NewReader(loaderTestData)
	docs, _ := fetchDocs(reader)
	flags := make([]bool, len(docs))
	for idx, doc := range docs {
		flags[idx] = loader.shouldUseDoc(doc)
	}
	if !flags[0] || !flags[1] || flags[2] {
		t.Errorf("Result mismatch got %v flags, expected [true, true, false]", flags)
	}
}
