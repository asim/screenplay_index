package app

import (
	"github.com/dchest/captcha"
	"log"
	"net/http"
)

func Logger(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header.Get("X-Forwarded-For"), r.RemoteAddr, r.Method, r.RequestURI, r.Referer(), r.UserAgent())
}

func Run(host string) {
	// admin handlers
	//http.HandleFunc("/_add", adderHandler)
	http.HandleFunc("/_pending", pendingHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/scripts", scriptsHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/s/", shortHandler)
	log.Println("Starting listening on", host)
	log.Fatal(http.ListenAndServe(host, nil))

}
