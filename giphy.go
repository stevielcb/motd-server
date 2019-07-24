package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "strconv"
)

var (
  apiKey string
)

func init() {
  dat, err := ioutil.ReadFile(giphyKeyFile)
  if err != nil {
    panic(err)
  }

  apiKey = string(dat)
}

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

  body, err := ioutil.ReadAll(resp.Body)
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
