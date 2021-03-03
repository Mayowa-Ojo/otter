package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	yaml "github.com/goccy/go-yaml"
	"github.com/theckman/yacspin"
)

// PERMISSION - file mode
const PERMISSION os.FileMode = 0777

//CONFIG_PATH - otter config location
const CONFIG_PATH string = "/.config/otter"

// OpenURLLink - opens specified url in default os browser
func OpenURLLink(url string) error {
	var cmd string // os-specific command to launch browser
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "cmd"
	}

	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}

// TokenPair - pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// GetAuthTokens - fetch auth tokens from local conf
func GetAuthTokens() (*TokenPair, error) {
	var tokens *TokenPair
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	byt, err := ioutil.ReadFile(homeDir + CONFIG_PATH + "/.keys")

	if err != nil {
		return nil, err
	}

	content := string(byt)
	lines := strings.Split(content, "\n")

	if !strings.HasPrefix(lines[0], "access_token=") {
		return nil, errors.New("no auth token found")
	}

	if !strings.HasPrefix(lines[1], "refresh_token=") {
		return nil, errors.New("no auth token found")
	}

	accessToken := strings.Split(lines[0], "=")[1]
	refreshToken := strings.Split(lines[1], "=")[1]

	if isTokenValid := VerifyAuthToken(accessToken); isTokenValid {
		tokens = &TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		return tokens, nil
	}

	tokens, err = UpdateAuthorization(refreshToken)
	if err != nil {
		return nil, errors.New("authorization failed")
	}

	return tokens, nil
}

// VerifyAuthToken - check if provided token is valid
// [token] - access token
func VerifyAuthToken(token string) bool {
	url := "https://api.heroku.com/apps"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)

	if err != nil {
		return false
	}

	return resp.StatusCode == 200
}

// PersistAuthorization - save auth tokens to user system
// [at] - access token
// [rt] - refresh token
func PersistAuthorization(at string, rt string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	content := fmt.Sprintf("access_token=%s\nrefresh_token=%s", at, rt)

	err = ioutil.WriteFile(homeDir+CONFIG_PATH+"/.keys", []byte(content), os.FileMode(PERMISSION))

	if !os.IsNotExist(err) {
		return err
	}

	// create config path
	err = os.Mkdir(homeDir+CONFIG_PATH, os.FileMode(PERMISSION))

	err = ioutil.WriteFile(homeDir+CONFIG_PATH+"/.keys", []byte(content), os.FileMode(PERMISSION))

	return err
}

// UpdateAuthorization - get new access_token from refresh token
func UpdateAuthorization(refreshToken string) (*TokenPair, error) {
	var tokens TokenPair
	// uri := fmt.Sprintf("http://localhost:5000/auth/refresh")
	uri := fmt.Sprintf("https://otter-api-server.herokuapp.com/auth/refresh")

	body, err := json.Marshal(map[string]interface{}{
		"refresh_token": refreshToken,
	})

	b := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST", uri, b)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("authorization failed")
	}

	var data map[string]interface{}

	byt, err := io.ReadAll(resp.Body)

	if err := json.Unmarshal(byt, &data); err != nil {
		return nil, err
	}

	if _, ok := data["access_token"]; !ok {
		return nil, errors.New("authorization failed")
	}

	accessToken := data["access_token"].(string)
	refreshToken = data["refresh_token"].(string)

	if err := PersistAuthorization(accessToken, refreshToken); err != nil {
		return nil, err
	}

	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return &tokens, nil
}

// RevokeAuthorization - invalidate all tokens from user system
func RevokeAuthorization() error {
	homeDir, err := os.UserHomeDir()
	tokens, err := GetAuthTokens()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("https://api.heroku.com/authorizations/%s", tokens.AccessToken)
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		return errors.New("failed to revoke authorization")
	}

	err = ioutil.WriteFile(homeDir+CONFIG_PATH+"/.keys", []byte(""), os.FileMode(PERMISSION))

	return err
}

// ParseEnv - convert env file to map structure
// [path] - relative path to env file
func ParseEnv(path string) (map[string]string, error) {
	var out map[string]string

	byt, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(byt)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		kv := strings.Split(line, "=")
		key := kv[0]
		value := kv[1]

		out[key] = value
	}

	return out, nil
}

// ParseYAML - convert yaml file to map structure
// [path] - relative path to yaml file
func ParseYAML(path string) (map[string]string, error) {
	var out map[string]string

	if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, "yml") {
		return nil, errors.New("invalid path for yaml file")
	}

	byt, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(byt, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// ParseJSON - convert json file to map structure
// [path] - relative path to json file
func ParseJSON(path string) (map[string]string, error) {
	var out map[string]string

	if !strings.HasSuffix(path, ".json") {
		return nil, errors.New("invalid path for json file")
	}

	byt, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byt, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// LoadingSpinner - show loading spinner
func LoadingSpinner() (*yacspin.Spinner, error) {
	config := yacspin.Config{
		Frequency:         100 * time.Millisecond,
		CharSet:           yacspin.CharSets[69],
		Prefix:            "Please wait... ",
		Message:           "",
		StopCharacter:     "✓",
		StopFailMessage:   "",
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
		StopColors:        []string{"fgGreen"},
	}

	spinner, err := yacspin.New(config)
	if err != nil {
		return nil, err
	}

	return spinner, nil
}

// GenerateDataTable - create a structure table to display dataset
// [dataset] - data to be rendered
func GenerateDataTable(dataset []interface{}) (*simpletable.Table, error) {
	table := simpletable.New()

	switch dataset[0].(type) {
	// NOTE: not very satisfied with this case value. it's very generic and wille easily break.
	case map[string]interface{}:
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "#"},
				{Align: simpletable.AlignCenter, Text: "Key"},
				{Align: simpletable.AlignCenter, Text: "Value"},
			},
		}

		for i, v := range dataset {
			data := v.(map[string]interface{})
			r := []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", i)},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", data["key"])},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", data["value"])},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}

	default:
		return nil, errors.New("invalid type for dataset")
	}

	table.SetStyle(simpletable.StyleMarkdown)

	return table, nil
}
