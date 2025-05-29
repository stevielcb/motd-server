// giphy.go
//
// This file handles fetching random Giphy images to be used as MOTDs (messages of the day).
// It interacts with the Giphy API to fetch either an original or a downsized image based on size limits.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var (
	apiKey string
)

// init reads the Giphy API key from the file specified by giphyKeyFile
// and initializes the global apiKey variable.
func init() {
	dat, err := os.ReadFile(giphyKeyFile)
	if err != nil {
		panic(err)
	}

	apiKey = string(dat)
}

// randomGiphy fetches a random Giphy URL matching the given tag and rating.
//
// It selects an image and checks the size of the "original" version.
// If the original exceeds 10MB, it falls back to the downsized large version.
// Returns the URL of the selected Giphy image or an error if encountered.
func randomGiphy(tag string, rating string) (string, error) {
	url := fmt.Sprintf(
		"http://api.giphy.com/v1/gifs/random?api_key=%s&tag=%s&rating=%s",
		apiKey,
		url.QueryEscape(tag),
		rating,
	)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	data := result["data"].(map[string]interface{})
	if len(data) < 1 {
		return "", nil
	}

	images := data["images"].(map[string]interface{})
	original := images["original"].(map[string]interface{})
	downsized := images["downsized_large"].(map[string]interface{})

	sizeResp, err := http.Head(original["url"].(string))
	if err != nil {
		return "", err
	}

	if sizeResp.StatusCode != http.StatusOK {
		return "", err
	}

	size, _ := strconv.Atoi(sizeResp.Header.Get("Content-Length"))
	// If the original gif is larger than 10MB,
	// get the downsized image instead.
	if int64(size) > 10485760 {
		return downsized["url"].(string), nil
	}

	return original["url"].(string), nil
}
