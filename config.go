package main

import (
  "github.com/kelseyhightower/envconfig"
  "os"
  "strings"
)

var (
  cacheDir     string
  start        string
  end          string
  giphyKeyFile string
  c            Config
)

type Config struct {
  CacheDir         string            `split_words:"true"`
  CacheMaxFiles    int               `split_words:"true" default:"50"`
  GiphyApiKeyFile  string            `split_words:"true"`
  GiphyTags        map[string]string `split_words:"true"`
  DownloadInterval int               `split_words:"true" default:"10"`
  CleanupInterval  int               `split_words:"true" default:"60"`
  ListenHost       string            `split_words:"true" default:"localhost"`
  ListenPort       int               `split_words:"true" default:"4200"`
}

func init() {
  err := envconfig.Process("motd", &c)
  if err != nil {
    panic(err)
  }

  _, ok := os.LookupEnv("TERM")
  if ok {
    term := os.Getenv("TERM")
    if strings.HasPrefix(term, "screen") {
      start = "\033Ptmux;\033\033]"
      end = "\a\033\\"
    } else {
      start = "\033]"
      end = "\a"
    }
  } else {
    os.Exit(1)
  }

  home := os.Getenv("HOME")

  if c.GiphyApiKeyFile == "" {
    giphyKeyFile = home + "/.giphy-api"
  } else {
    giphyKeyFile = c.GiphyApiKeyFile
  }

  if c.CacheDir == "" {
    cacheDir = home + "/.motd"
  } else {
    cacheDir = c.CacheDir
  }

  os.Mkdir(cacheDir, 0700)
}
