package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     "86gkduq4srneks",
		ClientSecret: "5ZLK7n8nid9nZDvT",
		RedirectURL:  "https://localhost:9090/oauth2callback",
		Scopes:       []string{"r_basicprofile"},
		Endpoint:     linkedin.Endpoint,
	}
	oauthStateString = "thisshouldberandom"
)

const redirectURL = "https://192.168.1.70:4200/#!/profile-manager/new-profile"

func main() {
	http.HandleFunc("/login", handleFacebookLogin)
	http.HandleFunc("/oauth2callback", handleFacebookCallback)
	fmt.Print("Started running on http://localhost:9090\n")
	err := http.ListenAndServeTLS(":9090", "../mycert.pem", "../server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// const htmlIndex = `<html><body>
// Logged in with <a href="/login">facebook</a>
// </body></html>
// `

// func handleMain(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(htmlIndex))
// }

func handleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	Url, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	fmt.Printf(Url.Path)
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, redirectURL+"?fb_error=1", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		http.Redirect(w, r, redirectURL+"?fb_error=1", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
		http.Redirect(w, r, redirectURL+"?fb_error=1", http.StatusTemporaryRedirect)
		return
	}

	fmt.Printf("token: %s\n", token)
	fmt.Printf("parseResponseBody: %s\n", string(response))

	http.Redirect(w, r, redirectURL+"?fb_access_token="+token.AccessToken, http.StatusTemporaryRedirect)
}
