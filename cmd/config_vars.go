package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Mayowa-Ojo/otter/internal"
	// "strings"
)

const baseURI = "https://api.heroku.com"

// ConfigVar - environment variable key-value pair
type ConfigVar struct {
	key   string
	value string
}

// GetVariables - fetch all config vars for given app
// [app] - app name or id
func GetVariables(app, token string) (map[string]interface{}, error) {
	uri := fmt.Sprintf("%s/apps/%s/config-vars", baseURI, app)
	client := &http.Client{}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			return nil, errors.New("client is not authorized")
		}
		return nil, errors.New("error fetching resource")
	}

	var data map[string]interface{}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// UpsertVariable - add or update an existing variable
// [app] - app name or id
func UpsertVariable(app, token string, variable ConfigVar) error {
	uri := fmt.Sprintf("%s/apps/%s/config-vars", baseURI, app)
	client := &http.Client{}

	body, err := json.Marshal(map[string]string{
		variable.key: variable.value,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", uri, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			return errors.New("client is not authorized")
		}
		return errors.New("error adding resource")
	}

	return nil
}

// UpsertVariables - add or update existing variables.
// [path] - relative file path
// [source] - can be a json, yaml or .env file
// [app] - app name or id
func UpsertVariables(app, token, path, source string) error {
	uri := fmt.Sprintf("%s/apps/%s/config-vars", baseURI, app)
	client := &http.Client{}
	var content map[string]string

	switch source {
	case "env":
		c, err := internal.ParseEnv(path)
		if err != nil {
			return err
		}

		content = c
	case "json":
		c, err := internal.ParseJSON(path)
		if err != nil {
			return err
		}

		content = c
	case "yaml":
		c, err := internal.ParseYAML(path)
		if err != nil {
			return err
		}

		content = c
	}

	body, err := json.Marshal(content)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", uri, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			return errors.New("client is not authorized")
		}
		return errors.New("error adding resource")
	}

	var data map[string]interface{}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	return nil
}

// RemoveVariable - remove an existing variable
// [app] - app name or id
// [token] - access token
// [key] - variable to be removed
func RemoveVariable(app, token, key string) error {
	uri := fmt.Sprintf("%s/apps/%s/config-vars", baseURI, app)
	client := &http.Client{}

	body, err := json.Marshal(map[string]interface{}{
		key: nil,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", uri, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			return errors.New("client is not authorized")
		}
		return errors.New("error adding resource")
	}

	return nil
}
