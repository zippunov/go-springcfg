package springcfg

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/thoas/go-funk"
)

// Loader loads configuration from multiple places
// - from local files application.yml, application-<profile>.yml
// - from spring cloud service if bootstrap yaml allow
type Loader struct {
	ActiveProfiles []string
	Dir            string
}

// LoadConfig returns fetched Data
func (l *Loader) LoadConfig() (Data, error) {
	rslt := Data{}
	bootstrap, err := l.loadBootsrap()
	if err != nil {
		return rslt, err
	}
	rslt, err = l.loadLocalConfig()
	if err != nil {
		return rslt, err
	}
	if bootstrap.GetBool("spring.cloud.config.enabled") {
		configClient := SpringCloudConfigClient{
			BaseURL:    bootstrap.GetString("spring.cloud.config.uri"),
			VaultToken: bootstrap.GetString("spring.cloud.config.vault_token"),
			Name:       bootstrap.GetString("spring.application.name"),
			Profiles:   l.ActiveProfiles,
			Label:      bootstrap.GetString("spring.cloud.config.label"),
		}
		cloudData, err := configClient.Fetch()
		if err == nil {
			rslt = rslt.Merge(cloudData)
		}
	}

	if bootstrap.Has("spring.config.location") {
		// load external yamls from the local paths
		regex := regexp.MustCompile(`^file:(\/\w.+\.ya?ml)$`)
		setting := bootstrap.GetString("spring.config.location")
		tokens := strings.Split(setting, ",")
		tokens = funk.FilterString(tokens, func(s string) bool {
			return regex.MatchString(s)
		})
		paths := make([]string, len(tokens))
		for idx, t := range tokens {
			paths[idx] = regex.FindStringSubmatch(t)[1]
		}
		extData, err := l.loadMerged(paths)
		if err == nil {
			rslt = rslt.Merge(extData)
		}

	}

	return rslt, nil
}

func (l *Loader) shouldUseDoc(doc map[string]interface{}) bool {
	d := NewData(doc)
	field := d.GetString("profiles")
	if field == "" {
		return true
	}
	profiles := strings.Split(field, ",")
	includedProfiles := map[string]interface{}{}
	excludedProfiles := map[string]interface{}{}
	flag := struct{}{}
	for _, profile := range profiles {
		if strings.HasPrefix(profile, "!") {
			cleanProfile := strings.TrimSpace(profile[1:])
			if cleanProfile != "" {
				excludedProfiles[cleanProfile] = flag
			}
		} else {
			cleanProfile := strings.TrimSpace(profile)
			if cleanProfile != "" {
				includedProfiles[cleanProfile] = flag
			}
		}
	}
	for _, profile := range l.ActiveProfiles {
		if _, ok := excludedProfiles[profile]; ok {
			return false
		}
	}
	for _, profile := range l.ActiveProfiles {
		if _, ok := includedProfiles[profile]; ok {
			return true
		}
	}
	return false
}

func (l *Loader) loadYaml(path string) (Data, error) {
	f, fileErr := os.Open(path)
	if fileErr != nil {
		return Data{}, fileErr
	}
	defer f.Close()
	docs, yamlErr := fetchDocs(bufio.NewReader(f))
	if yamlErr != nil {
		return Data{}, yamlErr
	}
	filteredDocs := docs[:0]
	for _, doc := range docs {
		if l.shouldUseDoc(doc) {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return NewData(filteredDocs...), nil
}

func (l *Loader) loadBootsrap() (Data, error) {
	bootloaderPath := filepath.Join(l.Dir, "bootstrap.yml")
	rslt, err := l.loadYaml(bootloaderPath)
	if err != nil && os.IsNotExist(err) {
		return NewData(), nil
	}
	return rslt, err
}

func (l *Loader) loadLocalConfig() (Data, error) {
	paths := []string{filepath.Join(l.Dir, "application.yml")}
	for _, profile := range l.ActiveProfiles {
		paths = append(paths, filepath.Join(l.Dir, "application-"+profile+".yml"))
	}
	return l.loadMerged(paths)
}

func (l *Loader) loadMerged(paths []string) (Data, error) {
	data := Data{}
	for _, p := range paths {
		nd, err := l.loadYaml(p)
		if err != nil {
			return data, err
		}
		data = data.Merge(nd)
	}
	return data, nil
}
