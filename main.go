package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
)

var (
	cookieName    = "jh-proxy-auth"
	hubApi        = os.Getenv("JUPYTERHUB_API_URL")
	apiToken      = os.Getenv("JUPYTERHUB_API_TOKEN")
	clientId      = os.Getenv("JUPYTERHUB_CLIENT_ID")
	callbackUrl   = os.Getenv("JUPYTERHUB_OAUTH_CALLBACK_URL")
	servicePrefix = os.Getenv("JUPYTERHUB_SERVICE_PREFIX")
	jhUser        = os.Getenv("JUPYTERHUB_USER")
	static        = "/static/desktop/"
	cookieSource  = securecookie.New(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))
	target        = flag.String("target", "http://127.0.0.1:8080", "the target host/port")
	port          = flag.String("port", "8888", "the port to serve on")
)

type JHOAuthHandler struct {
	wrappedHandler http.Handler
}

func validateCookie(r *http.Request) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Println(err)
		return false
	}

	var value string
	if err = cookieSource.Decode(cookieName, cookie.Value, &value); err != nil {
		log.Println(err)
		return false
	}

	req, err := http.NewRequest("GET", hubApi+"/authorizations/token/"+value, nil)
	req.Header.Add("Authorization", "Bearer "+apiToken)
	if res, err := http.DefaultClient.Do(req); err == nil {
		payload, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			return false
		}
		var resp map[string]interface{}
		err = json.Unmarshal(payload, &resp)
		if err != nil {
			log.Println(err)
			return false
		}
		if resp["name"] == jhUser {
			return true
		}
	} else {
		log.Println(err)
	}
	return false
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if code := r.URL.Query().Get("code"); code == "" {
		w.WriteHeader(400)
	} else {
		res, err := http.PostForm(hubApi+"/oauth2/token", url.Values{
			"client_id":     {clientId},
			"client_secret": {apiToken},
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {callbackUrl},
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		if res.StatusCode == 200 {
			var resp map[string]interface{}
			payload, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}
			err = json.Unmarshal(payload, &resp)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}
			if encoded, err := cookieSource.Encode(cookieName, resp["access_token"]); err == nil {
				cookie := &http.Cookie{
					Name:     cookieName,
					Value:    encoded,
					Path:     servicePrefix,
					Secure:   true,
					HttpOnly: true,
				}
				http.SetCookie(w, cookie)
				cookie.Path = strings.Replace(servicePrefix, "@", "%40", -1)
				http.SetCookie(w, cookie)
				http.Redirect(w, r, servicePrefix, http.StatusFound)
			} else {
				log.Println(err)
				w.WriteHeader(500)
			}
		}
	}
}

func (ah JHOAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	if validateCookie(r) {
		ah.wrappedHandler.ServeHTTP(w, r)
	} else if r.URL.Path == callbackUrl {
		handleCallback(w, r)
	} else {
		cookie := http.Cookie{Name: "jh-proxy-auth-state", Value: r.URL.String(), MaxAge: 600, Path: servicePrefix}
		http.SetCookie(w, &cookie)
		params := fmt.Sprintf("?client_id=%s&redirect_uri=%s&response_type=code&state=", url.QueryEscape(clientId), url.QueryEscape(callbackUrl))
		http.Redirect(w, r, "/hub/api/oauth2/authorize"+params, http.StatusFound)
	}
}

func newPathTrimmingReverseProxy(target *url.URL) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host

			// Static files need to maintain the complete request path
			if !strings.Contains(req.URL.Path, static) {
				req.URL.Path = strings.TrimPrefix(req.URL.Path, strings.TrimSuffix(servicePrefix, "/"))
				req.URL.RawPath = strings.TrimPrefix(req.URL.RawPath, strings.TrimSuffix(servicePrefix, "/"))
			}
			if _, ok := req.Header["User-Agent"]; !ok {
				req.Header.Set("User-Agent", "") // explicitly disable User-Agent so it's not set to default value
			}
		},
	}
}

func main() {
	flag.Parse()
	backend, err := url.Parse(*target)
	if err != nil {
		log.Fatalln(err)
	}
	handler := JHOAuthHandler{
		wrappedHandler: newPathTrimmingReverseProxy(backend),
	}

	// wait until target is reachable
	for {
		res, err := http.Get(backend.String())
		if err == nil && res.StatusCode == 200 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	err = http.ListenAndServe(":"+*port, handler)
	if err != nil {
		log.Fatalln(err)
	}
}
