package redirects

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func jsonRedirects(jsonPath string) (Redirects, error) {
	contents, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read json file: %v", err)
	}

	redirects := map[string]string{}
	err = json.Unmarshal(contents, &redirects)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse json file: %v", err)
	}

	return redirects, nil
}
