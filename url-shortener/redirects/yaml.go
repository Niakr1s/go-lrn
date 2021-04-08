package redirects

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func yamlRedirects(yamlPath string) (Redirects, error) {
	contents, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read yaml file: %v", err)
	}

	redirects := map[string]string{}
	err = yaml.Unmarshal(contents, &redirects)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse yaml file: %v", err)
	}

	return redirects, nil
}
