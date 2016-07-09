package app

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/hoisie/mustache"
	elastigo "github.com/mattbaird/elastigo/lib"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"
)

var (
	templateCache = make(map[string]*mustache.Template)
	alphanum      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	pageLimit     = 20
)

type script struct {
	Id                      float64
	Meta, Title, Url, Short string
}

func (s *script) ID() string {
	return s.Short
}

func (s script) Domain() string {
	u, err := url.Parse(s.Url)
	if err != nil {
		return ""
	}
	return u.Host
}

func (s script) Uri() string {
	u, _ := url.QueryUnescape(s.Url)
	return url.QueryEscape(u)
}

func addScript(title, uri string) error {
	s := map[string]interface{}{
		"id":    float64(time.Now().Unix()),
		"meta":  "",
		"title": title,
		"url":   uri,
		"short": shorten(uri),
	}

	conn := elastigo.NewConn()

	rsp, err := conn.Index("scripts", "script", "", s, nil)
	if err != nil {
		log.Println("error indexing:", err)
		return err
	}

	log.Println("indexed item, response (id, type):", s, rsp.Id, rsp.Type)
	return nil
}

func alert(msg string) map[string]string {
	return map[string]string{
		"alert": msg,
	}
}

func urlExists(url string) bool {
	b, err := json.Marshal(url)
	if err != nil {
		log.Println(err)
		return true
	}

	conn := elastigo.NewConn()
	qs := elastigo.NewQueryString("url", string(b))
	q := elastigo.Query().Qs(&qs)
	out, err := elastigo.Search("scripts").Type("script").Query(q).Size("1").Result(conn)
	if err != nil {
		log.Println(err)
		return true
	}

	if out.Hits.Total > 0 {
		return true
	}

	return false
}

func shortExists(id string) bool {
	conn := elastigo.NewConn()
	out, err := elastigo.Search("scripts").Type("script").Search("short:" + id).Size("1").Result(conn)
	if err != nil {
		log.Println(err)
		return true
	}

	if out.Hits.Total > 0 {
		return true
	}

	return false
}

func shorten(url string) string {
	bytes := make([]byte, 10)
	for {
		rand.Read(bytes)
		for i, b := range bytes {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		}
		id := string(bytes)
		if !shortExists(id) {
			return id
		}
	}

	return "shortend"
}

func getIp(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); len(ip) > 0 {
		return ip
	}

	return r.RemoteAddr
}

func getPageOffset(vars url.Values, limit int) (int, int) {
	page, err := strconv.Atoi(vars.Get("page"))
	if err != nil {
		page = 1
	}

	if page > pageLimit {
		page = pageLimit
	}

	next := page - 1
	if page == 1 {
		next = 0
	}

	offset := next * limit
	return page, offset
}

func getPager(u *url.URL, page, limit, items int) map[string]string {
	pager := make(map[string]string)

	if page == 0 || page == 1 {
		pager["previousPage"] = "#"
		pager["previousState"] = "disabled"
	} else {
		prev := u
		vars := prev.Query()
		vars.Set("page", strconv.Itoa(page-1))
		prev.RawQuery = vars.Encode()
		pager["previousPage"] = prev.RequestURI()
	}

	if items < limit {
		pager["nextPage"] = "#"
		pager["nextState"] = "disabled"
	} else {
		next := u
		vars := next.Query()
		vars.Set("page", strconv.Itoa(page+1))
		next.RawQuery = vars.Encode()
		pager["nextPage"] = next.RequestURI()
	}

	return pager
}

func genTmpl(view string, data ...interface{}) string {
	tmpl, tprs := templateCache[view]
	lyot, lprs := templateCache["static/views/layout.m"]

	if !tprs {
		var err error
		tmpl, err = mustache.ParseFile(view)
		if err != nil {
			log.Println(err)
			return ""
		}
		templateCache[view] = tmpl
	}

	if !lprs {
		var err error
		lyot, err = mustache.ParseFile("static/views/layout.m")
		if err != nil {
			log.Println(err)
			return ""
		}
		templateCache["static/views/layout.m"] = lyot
	}

	return tmpl.RenderInLayout(lyot, data...)
}

func render(w http.ResponseWriter, r *http.Request, data interface{}, view string) {
	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(data)
		w.Write(b)
		return
	}

	viewPath := filepath.Join("static/views", view+".m")
	fmt.Fprintf(w, "%s", genTmpl(viewPath, data))
}
