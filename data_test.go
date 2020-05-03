package springcfg

import (
	"os"
	"strings"
	"testing"
)

var dataTestData = `
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
combine.prop: ${spring.application.name}:${spring.cloud.config.label}:${TEST_ENV_VAL1}
`

func TestConfigDataReplacements(t *testing.T) {
	os.Setenv("CONFIG_SERVER_VAULT_TOKEN", "xxxxxx")
	os.Setenv("TEST_ENV_VAL1", "amasiness")
	reader := strings.NewReader(dataTestData)
	docs, _ := fetchDocs(reader)
	data := NewData(docs...)
	if data.GetString("combine.prop") != "api-chat:development:amasiness" {
		t.Errorf("Result mismatch got \"%v\" , expected \"api-chat:development:amasiness\"", data.GetString("combine.prop"))
	}
	if data.GetString("spring.cloud.config.vault_token") != "xxxxxx" {
		t.Errorf("Result mismatch got \"%v\" , expected \"xxxxxx\"", data.GetString("spring.cloud.config.vault_token"))
	}
	if data.GetString("spring.cloud.config.uri") != "http://infra-config-server.service.consul:8888" {
		t.Errorf("Result mismatch got \"%v\" , expected \"http://infra-config-server.service.consul:8888\"", data.GetString("spring.cloud.config.uri"))
	}
}
