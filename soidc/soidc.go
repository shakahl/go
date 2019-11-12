package soidc

import (
	"net/http"
	"strings"

	log "github.com/kataras/golog"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"github.com/syncfuture/go/security"
	"golang.org/x/oauth2"
)

const (
	SESS_ID       = "ID"
	SESS_USERNAME = "Username"
	SESS_EMAIL    = "Email"
	SESS_ROLES    = "Roles"
	SESS_LEVEL    = "Level"
	SESS_STATUS   = "Status"
	SESS_STATE    = "State"
	COKI_TOKEN    = ".ART"
	COKI_SESSION  = ".USS"
)

type IOIDCClient interface {
	HandleAuthentication(context.Context)
	HandleSignInCallback(context.Context)
	HandleSignOut(context.Context)
	HandleSignOutCallback(context.Context)
	NewHttpClient(context.Context) (*http.Client, error)
	GetToken(ctx context.Context) (*oauth2.Token, error)
	SaveToken(ctx context.Context, token *oauth2.Token) error
}

type OIDCConfig struct {
	ProjectName       string
	ClientID          string
	ClientSecret      string
	ProviderUrl       string
	SignInCallbackURL string
	AccessDeniedURL   string
	Scopes            []string
}

type ClientOptions struct {
	ProviderURL        string
	ClientID           string
	ClientSecret       string
	SignInCallbackURL  string
	SignOutCallbackURL string
	AccessDeniedURL    string
	Sess_ID            string
	Sess_Username      string
	Sess_Email         string
	Sess_Roles         string
	Sess_Level         string
	Sess_Status        string
	Sess_State         string
	Coki_Token         string
	Coki_Session       string
	Scopes             []string
	Sessions           *sessions.Sessions
	SecureCookie       security.ISecureCookie
	PermissionAuditor  security.IPermissionAuditor
}

func getRoutes(handlerName string) (string, string, string) {
	array := strings.Split(handlerName, ".")
	return array[0], array[1], array[2]
}

func checkOptions(options *ClientOptions) {
	if options.ClientID == "" {
		log.Fatal("OIDCClient.Options.ClientID cannot be empty.")
	}
	if options.ClientSecret == "" {
		log.Fatal("OIDCClient.Options.ClientSecret cannot be empty.")
	}
	if len(options.Scopes) == 0 {
		log.Fatal("OIDCClient.Options.Scopes cannot be empty")
	}
	if options.ProviderURL == "" {
		log.Fatal("OIDCClient.Options.ProviderUrl cannot be empty.")
	}
	if options.SignInCallbackURL == "" {
		log.Fatal("OIDCClient.Options.SignInCallbackURL cannot be empty.")
	}
	if options.SignOutCallbackURL == "" {
		log.Fatal("OIDCClient.Options.SignOutCallbackURL cannot be empty.")
	}

	if options.Sessions == nil {
		log.Fatal("OIDCClient.Options.Sessions cannot be nil")
	}
	if options.PermissionAuditor == nil {
		log.Fatal("OIDCClient.Options.PermissionAuditor cannot be nil")
	}
	if options.SecureCookie == nil {
		log.Fatal("OIDCClient.Options.SecureCookie cannot be nil")
	}

	if options.Coki_Token == "" {
		options.Coki_Token = COKI_TOKEN
	}
	if options.Coki_Session == "" {
		options.Coki_Session = COKI_SESSION
	}

	if options.Sess_ID == "" {
		options.Sess_ID = SESS_ID
	}
	if options.Sess_Username == "" {
		options.Sess_Username = SESS_USERNAME
	}
	if options.Sess_Email == "" {
		options.Sess_Email = SESS_EMAIL
	}
	if options.Sess_Roles == "" {
		options.Sess_Roles = SESS_ROLES
	}
	if options.Sess_Level == "" {
		options.Sess_Level = SESS_LEVEL
	}
	if options.Sess_Status == "" {
		options.Sess_Status = SESS_STATUS
	}

	if options.Sess_State == "" {
		options.Sess_State = SESS_STATE
	}

	if options.AccessDeniedURL == "" {
		options.AccessDeniedURL = "/"
	}
}
