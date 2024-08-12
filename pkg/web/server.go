package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/sessions"
	"github.com/rusik69/iamrotator/pkg/aws"
	"github.com/rusik69/iamrotator/pkg/config"
	"github.com/rusik69/iamrotator/pkg/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var cfg config.Config

var store = sessions.NewCookieStore([]byte(cfg.Web.SSOCookieStoreKey))

var (
	cachedKeys       []types.AWSAccessKey
	cacheMutex       sync.Mutex
	cacheExpiration  time.Time
	cacheDuration    = 5 * time.Minute
	oauthConfig      *oauth2.Config
	oauthStateString string
)

// Listen starts the web server.
func Listen(configPath string) error {
	cf, err := config.Load(configPath)
	if err != nil {
		return err
	}
	cfg = cf
	oauthConfig = &oauth2.Config{
		RedirectURL:  cfg.Web.SSOCallbackURL,
		ClientID:     cfg.Web.SSOClientID,
		ClientSecret: cfg.Web.SSOClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = cfg.Web.SSOStateString
	store.MaxAge(3600 * 24) // 1 day
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	http.HandleFunc("/", serveHomePage)
	http.HandleFunc("/listkeys", isAuthenticated(serveListKeys))
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/auth/callback", handleGoogleCallback)
	logrus.Infof("Listening on %s", cfg.Web.ListenAddr)
	http.ListenAndServe(cfg.Web.ListenAddr, nil)
	return nil
}

// Middleware to check if user is authenticated
func isAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		session, err := store.Get(r, "session_token")
		if err != nil || session.Values["authenticated"] != true {
			logrus.Error("User is not authenticated")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// serveHomePage serves the home page.
func serveHomePage(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Serving home page")
	html, err := os.ReadFile("front/index.html")
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Could not read index.html", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// serveListKeys serves the list keys page.
func serveListKeys(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Serving list keys")
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if time.Now().After(cacheExpiration) {
		sess, err := aws.CreateSession(cfg.AWS)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		keys, err := aws.ListAccessKeys(sess, cfg.AWS)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cachedKeys = keys
		cacheExpiration = time.Now().Add(cacheDuration)
	}
	json.NewEncoder(w).Encode(cachedKeys)
}

// handleGoogleLogin initiates the Google login process.
func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGoogleCallback handles the callback from Google.
func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	userInfo := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if email, ok := userInfo["email"]; ok {
		found := false
		for _, allowedEmail := range cfg.Web.SSOAllowedEmails {
			if email == allowedEmail {
				found = true
				break
			}
		}
		if !found {
			logrus.Errorf("User %s not allowed", email)
			http.Error(w, "User not allowed", http.StatusForbidden)
			return
		}
	}
	sessionToken, err := generateRandomToken(32)
	if err != nil {
		http.Error(w, "Failed to generate session token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set session token
	session, _ := store.Get(r, "session-name")
	session.Values["session_token"] = sessionToken
	session.Save(r, w)
}

func generateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
