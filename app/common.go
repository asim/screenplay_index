package app

import (
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/mattbaird/elastigo/search"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
)

var (
	templateCache = make(map[string]*mustache.Template)
	alphanum      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

type script struct {
	Id                int64
	Title, Url, Short string
}

func alert(msg string) map[string]string {
	return map[string]string{
		"alert": msg,
	}
}

func exists(url string) bool {
	id := shorten(url)
	out, err := search.Search("scripts").Type("script").Search("short:" + id).Size("1").Result()
	if err != nil {
		return false
	}

	if out.Hits.Len() > 0 {
		return true
	}

	return false
}

func shorten(url string) string {
	bytes := make([]byte, 10)
	for i, b := range url {
		bytes[i%10] = alphanum[byte(b)%byte(len(alphanum))]
	}
	return string(bytes)
}

func getPageOffset(vars url.Values, limit int) (int, int) {
	page, err := strconv.Atoi(vars.Get("page"))
	if err != nil {
		page = 1
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

func render(w http.ResponseWriter, data interface{}, view string) {
	viewPath := filepath.Join("static/views", view+".m")
	fmt.Fprintf(w, "%s", genTmpl(viewPath, data))
}
