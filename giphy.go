package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
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

  return original["url"].(string), nil
}
