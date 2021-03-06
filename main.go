package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Error checking
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Config of OAuth2.0
var (
	Oauth_Config = &oauth2.Config{
		RedirectURL:  "http://localhost/api/callback",
		ClientID:     "266082700797-paa0aout6ieig192jolv0cv0oklm8qmh.apps.googleusercontent.com",
		ClientSecret: "HVTLT04u6Egv4LKRUN9yVMOP",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
)

// Random string generator
func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Handle Webpage
func main() {
	http.HandleFunc("/", hand_home)
	http.HandleFunc("/login", hand_login)
	http.HandleFunc("/callback", hand_callback)
	http.ListenAndServe(":8888", nil)
}

// html show at home
// Has login with Google feature
func hand_home(w http.ResponseWriter, r *http.Request) {
	home_html := `<html>
	<body>
		<a href="/api/login"> Click here to login with google </a>
	</body> <html>`
	fmt.Fprintf(w, home_html)
}

// Redirect login request to OAuth 2.0
func hand_login(w http.ResponseWriter, r *http.Request) {
	// State is a 32 bit random string that prevents CSRF attacks
	state_string := RandStringRunes(32)
	url := Oauth_Config.AuthCodeURL(state_string)

	// Login using OAuth 2.0
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Handle Callback from OAuth 2.0
func hand_callback(w http.ResponseWriter, r *http.Request) {
	// Get token
	code := r.FormValue("code")
	token, err := Oauth_Config.Exchange(oauth2.NoContext, code)
	check(err)

	// Exchange token for result
	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	check(err)
	defer res.Body.Close()

	// Print user info
	data, err := ioutil.ReadAll(res.Body)
	check(err)
	fmt.Fprintf(w, "User info : %s\n", data)
}
