package config

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	State           string           `json:"state"`
	PropertySources []propertySource `json:"propertySources"`
}

type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

// SpringCloudConfigClient fetches config from the Spring Cloud Config server
type SpringCloudConfigClient struct {
	BaseURL    string
	VaultToken string
	Name       string
	Profiles   []string
	Label      string
}

func (c *SpringCloudConfigClient) buildProfilesString() string {
	if c.Profiles == nil {
		return "default"
	}
	return strings.Join(c.Profiles, ",")
}

func (c *SpringCloudConfigClient) buildURL() string {
	parts := []string{
		c.BaseURL,
		url.PathEscape(c.Name),
		url.PathEscape(c.buildProfilesString()),
		url.PathEscape(c.Label),
	}
	return strings.Join(parts, "/")
}

// Fetch returns instance of the Data fetched from the spring cloud server
func (c *SpringCloudConfigClient) Fetch() (Data, error) {
	url := c.buildURL()
	client := &http.Client{}

	req, urlErr := http.NewRequest("GET", url, nil)
	if urlErr != nil {
		return Data{}, urlErr
	}
	req.Header.Add("Accept", `application/json, application/*+json`)
	if c.VaultToken != "" {
		req.Header.Add("X-Config-Token", c.VaultToken)
	}

	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return Data{}, fetchErr
	}
	defer resp.Body.Close()

	configBody := &springCloudConfig{}
	jsonErr := json.NewDecoder(resp.Body).Decode(configBody)
	if jsonErr != nil {
		return Data{}, jsonErr
	}
	docs := []map[string]interface{}{}
	for _, source := range configBody.PropertySources {
		docs = append(docs, source.Source)
	}
	return NewData(docs...), nil
}
