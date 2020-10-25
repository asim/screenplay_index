package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

type Result struct {
	Results []Script `json:"results"`
}

type Script struct {
	Short string `json:"short"`
	Title string `json":"title"`
}

var (
	seen = make(map[string]bool)
)

func getScripts(uri string) ([]Script, error) {
	r, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	rsp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var results Result
	if err := json.Unmarshal(b, &results); err != nil {
		return nil, err
	}

	return results.Results, nil
}

func tweet(uri, msg string) {
	scripts, err := getScripts(uri)
	if err != nil {
		fmt.Println("error getting latest", err)
		return
	}

	if len(scripts) == 0 {
		return
	}

	for i := 0; i < len(scripts); i++ {
		script := fmt.Sprintf("http://screenplays.com/s/%s", scripts[i].Short)
		title := scripts[i].Title
		tweetr:= fmt.Sprintf(msg, title, script)
		if seen[tweetr] {
			continue
		}
		seen[tweetr] = true
		api.PostTweet(tweetr, url.Values{})
		return
	}
}

func main() {
	anaconda.SetConsumerKey("")
	anaconda.SetConsumerSecret("")
	api := anaconda.NewTwitterApi("", "")

	// tickers
	clear := time.NewTicker(time.Hour * 72)
	latest := time.NewTicker(time.Hour * 24)
	random := time.NewTicker(time.Hour * 3)
	trending := time.NewTicker(time.Hour * 12)

	for {
		select {
		case <-clear.C:
			seen = make(map[string]bool)
		case <-latest.C:
			tweet("http://screenplays.com/scripts", `Latest on #screenplays: "%s". Read the script %s #screenwriting`)
		case <-random.C:
			tweet("http://screenplays.com/random", `Random reading on #screenplays: "%s". Read the script %s #screenwriting`)
		case <-trending.C:
			tweet("http://screenplays.com/trending", `Trending on #screenplays: "%s". Read the script %s #screenwriting`)
		}
	}
}
