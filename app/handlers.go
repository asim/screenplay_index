package app

import (
	"encoding/json"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/mattbaird/elastigo/core"
	"github.com/mattbaird/elastigo/search"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	itemLimit = 20
	itemSize  = "20"
)

func adderHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	title := strings.TrimSpace(r.FormValue("title"))
	uri := strings.TrimSpace(r.FormValue("url"))

	if title == "" || uri == "" {
		return
	}

	if _, err := url.Parse(uri); err != nil || !strings.HasSuffix(uri, ".pdf") {
		return
	}

	if t := urlExists(uri); t {
		log.Println("Exists", uri, t)
		return
	}

	s := map[string]interface{}{
		"id":    time.Now().Unix(),
		"title": title,
		"url":   uri,
		"short": shorten(uri),
	}

	rsp, err := core.Index(true, "scripts", "script", "", s)
	if err != nil {
		log.Println("error indexing:", err)
	} else {
		log.Println("indexed item %v response (id, type):", s, rsp.Id, rsp.Type)
	}

	return
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	if r.Method == "GET" {
		d := map[string]string{
			"captcha": captcha.New(),
		}

		render(w, d, "add")
		return
	}

	if r.Method != "POST" {
		fmt.Fprint(w, "Invalid HTTP Method")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	uri := strings.TrimSpace(r.FormValue("url"))
	digits := r.FormValue("captcha")
	id := r.FormValue("_captchaId")

	d := map[string]string{
		"captcha": captcha.New(),
	}

	if title == "" || uri == "" || digits == "" || id == "" {
		d["alert"] = "Fill in form"
		render(w, d, "add")
		return
	}

	if _, err := url.Parse(uri); err != nil || !strings.HasSuffix(uri, ".pdf") {
		d["alert"] = "Invalid URL (.pdf only)"
		render(w, d, "add")
		return
	}

	if t := urlExists(uri); t {
		log.Println("Exists", uri, t)
		d["alert"] = "Script already exists"
		render(w, d, "add")
		return
	}

	if !captcha.VerifyString(id, digits) {
		d["alert"] = "Invalid CAPTCHA input"
		render(w, d, "add")
		return
	}

	s := map[string]interface{}{
		"id":    time.Now().Unix(),
		"title": title,
		"url":   uri,
		"short": shorten(uri),
	}

	rsp, err := core.Index(true, "scripts", "script", "", s)
	if err != nil {
		log.Println("error indexing:", err)
	} else {
		log.Println("indexed item %v response (id, type):", s, rsp.Id, rsp.Type)
	}

	render(w, alert("Successfully added script"), "index")
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	render(w, nil, "index")
	return
}

func scriptsHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	vars := r.URL.Query()
	limit := itemLimit
	size := itemSize

	page, offset := getPageOffset(vars, limit)
	from := fmt.Sprintf("%d", offset)

	sort := search.Sort("id").Desc()
	out, err := search.Search("scripts").Type("script").From(from).Size(size).Sort(sort).Result()
	if err != nil {
		log.Println("Error:", err)
		render(w, alert("Problem retrieving scripts"), "scripts")
		return
	}

	var scripts []script
	for _, hit := range out.Hits.Hits {
		var s script
		err := json.Unmarshal(hit.Source, &s)
		if err != nil {
			log.Println(err)
			continue
		}

		scripts = append(scripts, s)
	}

	d := map[string]interface{}{
		"results": scripts,
		"total":   out.Hits.Total,
	}

	if out.Hits.Len() == itemLimit || !(page < 2) {
		pager := getPager(r.URL, page, limit, out.Hits.Len())
		d["pager"] = &pager
	}

	render(w, d, "scripts")
	return
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	vars := r.URL.Query()
	limit := itemLimit
	size := itemSize

	page, offset := getPageOffset(vars, limit)
	from := fmt.Sprintf("%d", offset)

	q := vars.Get("q")
	if q == "" {
		log.Println("Error: Query length 0")
		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	if len(q) > 2048 {
		log.Println("Error: Query length", len(q))
		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	qe := url.QueryEscape(q)
	out, err := search.Search("scripts").Type("script").Search(qe).From(from).Size(size).Result()
	if err != nil {
		log.Println("Error:", err)
		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	var scripts []script
	for _, hit := range out.Hits.Hits {
		var s script
		err := json.Unmarshal(hit.Source, &s)
		if err != nil {
			log.Println(err)
			continue
		}

		scripts = append(scripts, s)
	}

	d := map[string]interface{}{
		"results": scripts,
		"query":   q,
		"total":   out.Hits.Total,
	}

	if out.Hits.Len() == itemLimit || !(page < 2) {
		pager := getPager(r.URL, page, limit, out.Hits.Len())
		d["pager"] = &pager
	}

	render(w, d, "search")
	return
}

func shortHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	id := strings.TrimPrefix(r.RequestURI, "/s/")
	out, err := search.Search("scripts").Type("script").Search("short:" + id).Size("1").Result()
	if err != nil {
		log.Println("Error:", err)
		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	if out.Hits.Len() != 1 {
		http.NotFound(w, r)
		return
	}

	var s script
	err = json.Unmarshal(out.Hits.Hits[0].Source, &s)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	// The redirect
	http.Redirect(w, r, s.Url, 302)
}
