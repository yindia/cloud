package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"go.akshayshah.org/connectauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

const (
	sessionName = "custom-auth-session-store"
)

// Config represents the configuration for the AuthServer
type Config struct {
	Provider     string // "google", "github", "facebook", or "okta"
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	SessionKey   string
}

// AuthServer represents the authorization server
type AuthServer struct {
	config       Config
	sessionStore *sessions.CookieStore
	state        string
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
	mu           sync.Mutex // Add a mutex for thread-safe operations
}

// NewAuthServer creates and initializes a new AuthServer
func NewAuthServer(config Config) (*AuthServer, error) {
	ctx := context.Background()

	var oauth2Config oauth2.Config
	var verifier *oidc.IDTokenVerifier
	fmt.Println("config.Provider", config.Provider)
	switch config.Provider {
	case "google":
		oauth2Config = oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		}
		provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
		if err != nil {
			return nil, fmt.Errorf("failed to get Google provider: %v", err)
		}
		verifier = provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	case "github":
		oauth2Config = oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Endpoint:     github.Endpoint,
			Scopes:       []string{"user:email"},
		}
		// GitHub doesn't support OIDC, so we'll need to handle token verification differently

	case "okta":
		provider, err := oidc.NewProvider(ctx, config.Issuer)
		if err != nil {
			return nil, fmt.Errorf("failed to get Okta provider: %v", err)
		}
		oauth2Config = oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		}
		verifier = provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}

	return &AuthServer{
		config:       config,
		sessionStore: sessions.NewCookieStore([]byte(config.SessionKey)),
		state:        generateState(),
		oauth2Config: oauth2Config,
		verifier:     verifier,
	}, nil
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// LoginHandler handles the login request
func (as *AuthServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	as.mu.Lock()
	authURL := as.oauth2Config.AuthCodeURL(as.state, oidc.Nonce(generateState()))
	as.mu.Unlock()

	http.Redirect(w, r, authURL, http.StatusFound)
}

// AuthCodeCallbackHandler handles the authorization code callback
func (as *AuthServer) AuthCodeCallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	as.mu.Lock()
	if r.URL.Query().Get("state") != as.state {
		as.mu.Unlock()
		http.Error(w, `{"error": "Invalid state"}`, http.StatusBadRequest)
		return
	}
	as.mu.Unlock()

	// Make sure the code was provided
	if r.URL.Query().Get("code") == "" {
		http.Error(w, `{"error": "The code was not returned or is not accessible"}`, http.StatusBadRequest)
		return
	}

	oauth2Token, err := as.oauth2Config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to exchange token: %v"}`, err), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	_, err = as.verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// http.Redirect(w, r, "/", http.StatusFound)
	w.Write([]byte(fmt.Sprintf(`{"access_token": "%s"}`, oauth2Token.AccessToken)))
}

// LogoutHandler handles the logout request
func (as *AuthServer) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	as.mu.Lock()
	session, err := as.sessionStore.Get(r, sessionName)
	as.mu.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	as.mu.Lock()
	delete(session.Values, "id_token")
	delete(session.Values, "access_token")
	err = session.Save(r, w)
	as.mu.Unlock()
	if err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// IsAuthenticated checks if the user is authenticated
func (as *AuthServer) IsAuthenticated(r *connectauth.Request) bool {
	// Check for bearer token in the Authorization header
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		// Verify the token
		_, err := as.verifier.Verify(context.Background(), token)
		if err == nil {
			return false
		}
		return true
	}
	return false
}

// getProfileData retrieves the user's profile data
func (as *AuthServer) getProfileData(r *http.Request) (map[string]interface{}, error) {
	as.mu.Lock()
	session, err := as.sessionStore.Get(r, "okta-hosted-login-session-store")
	as.mu.Unlock()
	if err != nil {
		return nil, err
	}

	accessToken, ok := session.Values["access_token"].(string)
	if !ok {
		return nil, fmt.Errorf("no access token found in session")
	}

	userInfo, err := as.oauth2Config.Client(r.Context(), &oauth2.Token{AccessToken: accessToken}).Get(os.Getenv("ISSUER") + "/v1/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get userinfo: %v", err)
	}
	defer userInfo.Body.Close()

	var profile map[string]interface{}
	if err := json.NewDecoder(userInfo.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo: %v", err)
	}

	return profile, nil
}
