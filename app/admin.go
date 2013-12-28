package app

import (
	"log"
	"io/ioutil"
	"sync"
	"bytes"
	"strings"
	"net/http"
	"net"
)

type admin struct {
	User, Pass string
}

type adminManager struct {
	mtx sync.RWMutex
	admins map[string]admin
}

var (
	adminFile = "admins"
	defaultAdminManager = newAdminManager()
)

func init() {
	defaultAdminManager.init()
}

func newAdminManager() *adminManager {
	return &adminManager{
		admins: make(map[string]admin),
	}
}

func (a *adminManager) init() {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	b, err := ioutil.ReadFile(adminFile)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(b)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}

		s := strings.Split(line, " ")
		if len(s) < 3 {
			continue
		}

		a.admins[s[0]] = admin{s[1], strings.TrimSuffix(s[2], "\n")}
	}
}

func (a *adminManager) auth(r *http.Request) bool {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		h, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return false
		}
		ip = h
	}

	user := r.FormValue("user")
	pass := r.FormValue("pass")

	a.mtx.RLock()
	defer a.mtx.RUnlock()

	// Authed ip?
	admin, ok := a.admins[ip]
	if !ok {
		return false
	}

	return user == admin.User && pass == admin.Pass
}

func (a *adminManager) authIP(r *http.Request) bool {
        ip := r.Header.Get("X-Forwarded-For")
        if ip == "" {
		h, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return false
		}
		ip = h
        }

	a.mtx.RLock()
        defer a.mtx.RUnlock()

        // Authed ip?
        _, exists := a.admins[ip]
	return exists
}

func (a *adminManager) get(r *http.Request) admin {
        ip := r.Header.Get("X-Forwarded-For")
        if ip == "" {
		h, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return admin{}
		}
		ip = h
        }

	a.mtx.RLock()
	defer a.mtx.RUnlock()
	admin := a.admins[ip]
	log.Println(admin)
	return admin
}
