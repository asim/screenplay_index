package app

import (
	"encoding/json"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/mattbaird/elastigo/search"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	itemLimit = 20
	itemSize  = "20"
)

func adderHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	if !defaultAdminManager.auth(r) {
		http.Redirect(w, r, r.Referer(), 403)
		return
	}

	if r.Method != "POST" {
		http.Redirect(w, r, r.Referer(), 500)
		return
	}

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

	addScript(title, uri)
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

	s := script{
		Id:    time.Now().Unix(),
		Meta:  "",
		Title: title,
		Url:   uri,
		Short: shorten(uri),
	}

	if err := defaultPendingManager.add(s); err != nil {
		d["alert"] = err.Error()
		render(w, d, "add")
		return
	}

	render(w, alert("Script will appear in index soon"), "index")
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	switch r.URL.Path {
	case "/", "/index.html", "/index.htm":
		render(w, nil, "index")
	default:
		w.WriteHeader(http.StatusNotFound)
		render(w, nil, "404")
	}
	return
}

func latestHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	sort := search.Sort("id").Desc()
	out, err := search.Search("scripts").Type("script").From("0").Size(itemSize).Sort(sort).Result()
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
	}

	render(w, d, "scripts")
	return
}

func pendingHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	if r.Method == "GET" {
		if !defaultAdminManager.authIP(r) {
			http.Redirect(w, r, r.Referer(), 403)
			return
		}

		scripts := defaultPendingManager.read()
		results := map[string]interface{}{
			"admin":   defaultAdminManager.get(r),
			"results": scripts,
			"total":   len(scripts),
		}
		render(w, results, "pending")
		return
	}

	if !defaultAdminManager.auth(r) {
		http.Redirect(w, r, r.Referer(), 403)
		return
	}

	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	url := r.FormValue("url")

	if r.Method == "POST" {
		if r.FormValue("_method") == "DELETE" {
			if err := defaultPendingManager.reject(id, url); err != nil {
				render(w, alert(err.Error()), "pending")
				return
			}

			http.Redirect(w, r, r.Referer(), 302)
			return
		}

		if err := defaultPendingManager.approve(id, url); err != nil {
			render(w, alert(err.Error()), "pending")
			return
		}

		http.Redirect(w, r, r.Referer(), 302)
		return
	}
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

	qs := search.QueryString{"", "", url.QueryEscape(q), "", "", []string{"title", "meta"}}
	//	qs := search.NewQueryString("title", url.QueryEscape(q))
	qe := search.Query().Qs(&qs)

	out, err := search.Search("scripts").Type("script").Query(qe).From(from).Size(size).Result()
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
		w.WriteHeader(http.StatusNotFound)
		render(w, nil, "404")
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
	http.Redirect(w, r, s.Url, 301)

	// Trending click event
	defaultTrendingManager.click(&s, r.Header.Get("X-Forwarded-For"))
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	vars := r.URL.Query()
	id := vars.Get("s")
	u := vars.Get("url")

	if len(id) == 0 || len(u) == 0 {
		w.WriteHeader(http.StatusNotFound)
		render(w, nil, "404")
		return
	}

	out, err := search.Search("scripts").Type("script").Search("short:" + id).Size("1").Result()
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusNotFound)
		render(w, nil, "404")
		return
	}

	if out.Hits.Len() != 1 {
		w.WriteHeader(http.StatusNotFound)
		render(w, nil, "404")
		return
	}

	var s script
	err = json.Unmarshal(out.Hits.Hits[0].Source, &s)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	if su, _ := url.QueryUnescape(s.Url); u != su {
		w.WriteHeader(http.StatusNotFound)
		render(w, nil, "404")
		return
	}

	// The redirect
	http.Redirect(w, r, s.Url, 301)

	// Trending click event
	defaultTrendingManager.click(&s, r.Header.Get("X-Forwarded-For"))
}

func trendingHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	trending := defaultTrendingManager.getTrending()

	var scripts []script

	if len(trending) < numRanked {
		sort := search.Sort("short").Desc()
		if out, err := search.Search("scripts").Type("script").From("0").Size("20").Sort(sort).Result(); err == nil {
			for _, hit := range out.Hits.Hits {
				var s script
				err := json.Unmarshal(hit.Source, &s)
				if err != nil {
					continue
				}

				scripts = append(scripts, s)
			}
		}
	} else {
		for _, script := range trending {
			scripts = append(scripts, *script)
		}
	}

	d := map[string]interface{}{
		"results": scripts,
	}

	render(w, d, "trending")
	return
}
