package app

import (
	"encoding/json"
	"fmt"
	"github.com/dchest/captcha"
	elastigo "github.com/mattbaird/elastigo/lib"
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

		render(w, r, d, "add")
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
		render(w, r, d, "add")
		return
	}

	if _, err := url.Parse(uri); err != nil || !strings.HasSuffix(uri, ".pdf") {
		d["alert"] = "Invalid URL (.pdf only)"
		render(w, r, d, "add")
		return
	}

	if t := urlExists(uri); t {
		log.Println("Exists", uri, t)
		d["alert"] = "Script already exists"
		render(w, r, d, "add")
		return
	}

	if !captcha.VerifyString(id, digits) {
		d["alert"] = "Invalid CAPTCHA input"
		render(w, r, d, "add")
		return
	}

	s := script{
		Id:    float64(time.Now().Unix()),
		Meta:  "",
		Title: title,
		Url:   uri,
		Short: shorten(uri),
	}

	if err := defaultPendingManager.add(s); err != nil {
		d["alert"] = err.Error()
		render(w, r, d, "add")
		return
	}

	render(w, r, alert("Script will appear in index soon"), "index")
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	switch r.URL.Path {
	case "/", "/index.html", "/index.htm":
		render(w, r, nil, "index")
	default:
		w.WriteHeader(http.StatusNotFound)
		render(w, r, nil, "404")
	}
	return
}

func latestHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	conn := elastigo.NewConn()
	sort := elastigo.Sort("id").Desc()
	out, err := elastigo.Search("scripts").Type("script").From("0").Size(itemSize).Sort(sort).Result(conn)
	if err != nil {
		log.Println("Error:", err)
		render(w, r, alert("Problem retrieving scripts"), "scripts")
		return
	}

	var scripts []script
	for _, hit := range out.Hits.Hits {
		var s script
		err := json.Unmarshal(*hit.Source, &s)
		if err != nil {
			log.Println(err)
			continue
		}

		scripts = append(scripts, s)
	}

	d := map[string]interface{}{
		"results": scripts,
	}

	render(w, r, d, "scripts")
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
		render(w, r, results, "pending")
		return
	}

	if !defaultAdminManager.auth(r) {
		http.Redirect(w, r, r.Referer(), 403)
		return
	}

	id, _ := strconv.ParseFloat(r.FormValue("id"), 64)
	url := r.FormValue("url")

	if r.Method == "POST" {
		if r.FormValue("_method") == "DELETE" {
			if err := defaultPendingManager.reject(id, url); err != nil {
				render(w, r, alert(err.Error()), "pending")
				return
			}

			http.Redirect(w, r, r.Referer(), 302)
			return
		}

		if err := defaultPendingManager.approve(id, url); err != nil {
			render(w, r, alert(err.Error()), "pending")
			return
		}

		http.Redirect(w, r, r.Referer(), 302)
		return
	}
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	defaultTrackingManager.track(getIp(r), r.URL.Path)
	count := defaultTrackingManager.get(getIp(r), r.URL.Path)

	if count > 100 {
		http.Redirect(w, r, "/", 302)
		return
	}

	args := map[string]interface{}{
		"query": map[string]interface{}{
			"function_score": map[string]interface{}{
				"query": map[string]interface{}{
					"match_all": map[string]interface{}{},
				},
				"random_score": map[string]interface{}{},
			},
		},
		"size": 1,
	}

	conn := elastigo.NewConn()

	out, err := conn.DoCommand("POST", "/scripts/script/_search", nil, args)
	if err != nil {
		log.Println("Error:", err)
		render(w, r, alert("Problem retrieving script"), "random")
		return
	}

	var retval elastigo.SearchResult

	if err := json.Unmarshal(out, &retval); err != nil {
		log.Println("Error:", err)
		render(w, r, alert("Problem retrieving script"), "random")
		return
	}

	var scripts []script
	for _, hit := range retval.Hits.Hits {
		var s script
		err := json.Unmarshal(*hit.Source, &s)
		if err != nil {
			log.Println(err)
			continue
		}

		scripts = append(scripts, s)
	}

	d := map[string]interface{}{
		"results": scripts,
		"total":   retval.Hits.Total,
	}

	render(w, r, d, "random")
	return
}

func scriptsHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	vars := r.URL.Query()
	limit := itemLimit
	size := itemSize

	page, offset := getPageOffset(vars, limit)
	from := fmt.Sprintf("%d", offset)
	conn := elastigo.NewConn()
	sort := elastigo.Sort("id").Desc()
	out, err := elastigo.Search("scripts").Type("script").From(from).Size(size).Sort(sort).Result(conn)
	if err != nil {
		log.Println("Error:", err)
		render(w, r, alert("Problem retrieving scripts"), "scripts")
		return
	}

	var scripts []script
	for _, hit := range out.Hits.Hits {
		var s script
		err := json.Unmarshal(*hit.Source, &s)
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

	render(w, r, d, "scripts")
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

	var nq string

	for _, r := range q {
		s := string(r)
		switch s {
		case "+", "-", "&&", "||", "!", "(", ")", "{", "}", "[", "]", "^", "\"", "~", "*", "?", ":", "\\", "/":
			nq += "\\"
		}
		nq += s
	}

	qs := elastigo.QueryString{
		DefaultOperator: "",
		DefaultField:    "",
		Query:           nq,
		Exists:          "",
		Missing:         "",
		Fields:          []string{"title", "meta"},
	}
	qe := elastigo.Query().Qs(&qs)
	conn := elastigo.NewConn()
	out, err := elastigo.Search("scripts").Type("script").Query(qe).From(from).Size(size).Result(conn)
	if err != nil {
		log.Println("Error:", err)
		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	var scripts []script
	for _, hit := range out.Hits.Hits {
		var s script
		err := json.Unmarshal(*hit.Source, &s)
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

	render(w, r, d, "search")
	return
}

func shortHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	conn := elastigo.NewConn()
	id := strings.TrimPrefix(r.RequestURI, "/s/")
	out, err := elastigo.Search("scripts").Type("script").Search("short:" + id).Size("1").Result(conn)
	if err != nil {
		log.Println("Error:", err)
		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	if out.Hits.Len() != 1 {
		w.WriteHeader(http.StatusNotFound)
		render(w, r, nil, "404")
		return
	}

	var s script
	err = json.Unmarshal(*out.Hits.Hits[0].Source, &s)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	// The redirect
	http.Redirect(w, r, s.Url, 301)

	// Trending click event
	defaultTrendingManager.click(&s, getIp(r))
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)
	vars := r.URL.Query()
	id := vars.Get("s")
	u := vars.Get("url")

	if len(id) == 0 || len(u) == 0 {
		w.WriteHeader(http.StatusNotFound)
		render(w, r, nil, "404")
		return
	}

	conn := elastigo.NewConn()
	out, err := elastigo.Search("scripts").Type("script").Search("short:" + id).Size("1").Result(conn)
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusNotFound)
		render(w, r, nil, "404")
		return
	}

	if out.Hits.Len() != 1 {
		w.WriteHeader(http.StatusNotFound)
		render(w, r, nil, "404")
		return
	}

	var s script
	err = json.Unmarshal(*out.Hits.Hits[0].Source, &s)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	if su, _ := url.QueryUnescape(s.Url); u != su {
		w.WriteHeader(http.StatusNotFound)
		render(w, r, nil, "404")
		return
	}

	// The redirect
	http.Redirect(w, r, s.Url, 301)

	// Trending click event
	defaultTrendingManager.click(&s, getIp(r))
}

func trendingHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	trending := defaultTrendingManager.getTrending()

	var scripts []script

	conn := elastigo.NewConn()

	if len(trending) < numRanked/2 {
		sort := elastigo.Sort("short").Desc()
		if out, err := elastigo.Search("scripts").Type("script").From("0").Size("20").Sort(sort).Result(conn); err == nil {
			for _, hit := range out.Hits.Hits {
				var s script
				err := json.Unmarshal(*hit.Source, &s)
				if err != nil {
					continue
				}

				scripts = append(scripts, s)
			}
		}
	} else {
		scripts = trending
	}

	d := map[string]interface{}{
		"results": scripts,
	}

	render(w, r, d, "trending")
	return
}
