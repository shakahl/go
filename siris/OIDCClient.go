package siris

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kataras/iris/v12/sessions"
	"github.com/syncfuture/go/rand"
	"github.com/syncfuture/go/security"

	oidc "github.com/coreos/go-oidc"
	"github.com/kataras/iris/v12/context"

	gocontext "context"

	"golang.org/x/oauth2"
)

type ClientOptions struct {
	ProviderUrl         string
	ClientID            string
	ClientSecret        string
	RedirectURL         string
	Session_ID          string
	Session_Username    string
	Session_Email       string
	Session_Roles       string
	Session_Level       string
	Session_Status      string
	Session_State       string
	Cookie_RefreshToken string
	Cookie_AccesToken   string
	Cookie_Session      string
	Scopes              []string
	Sessions            *sessions.Sessions
	PermissionAuditor   security.IPermissionAuditor
	SecureCookie        security.ISecureCookie
}

type IOIDCClient interface {
	HandleAuthentication(ctx context.Context)
	HandleSignInCallback(ctx context.Context)
	HandleSignOutCallback(ctx context.Context)
}

type defaultOIDCClient struct {
	Options      *ClientOptions
	OIDCProvider *oidc.Provider
	OAuth2Config *oauth2.Config
}

func NewOIDCClient(options *ClientOptions) IOIDCClient {
	checkOptions(options)

	x := new(defaultOIDCClient)
	x.Options = options

	ctx := gocontext.Background()
	var err error
	x.OIDCProvider, err = oidc.NewProvider(ctx, options.ProviderUrl)
	if err != nil {
		log.Fatal(err)
	}
	x.OAuth2Config = new(oauth2.Config)
	x.OAuth2Config.ClientID = options.ClientID
	x.OAuth2Config.ClientSecret = options.ClientSecret
	x.OAuth2Config.Endpoint = x.OIDCProvider.Endpoint()
	x.OAuth2Config.RedirectURL = options.RedirectURL
	x.OAuth2Config.Scopes = append(options.Scopes, oidc.ScopeOpenID)

	return x
}

func (x *defaultOIDCClient) HandleAuthentication(ctx context.Context) {
	session := x.Options.Sessions.Start(ctx)

	handlerName := ctx.GetCurrentRoute().MainHandlerName()
	area, controller, action := getRoutes(handlerName)

	// 判断请求是否允许访问
	userid := session.GetString(x.Options.Session_ID)
	if userid != "" {
		roles := session.GetInt64Default(x.Options.Session_Roles, 0)
		// 已登录
		allow := x.Options.PermissionAuditor.CheckRoute(area, controller, action, roles)
		if allow {
			// 有权限
			ctx.Next()
			return

		} else {
			// 没权限
			// Todo: 引导去提示页面
		}
	} else {
		// 未登录
		allow := x.Options.PermissionAuditor.CheckRoute(area, controller, action, 0)
		if allow {
			// 允许匿名
			ctx.Next()
			return
		}

		// 跳转去登录页面
		state := rand.String(32)
		session.Set(x.Options.Session_State, state)
		ctx.Redirect(x.OAuth2Config.AuthCodeURL(state), http.StatusFound)
	}
}

func (x *defaultOIDCClient) getRoutes(handlerName string) (string, string, string) {
	array := strings.Split(handlerName, ".")
	return array[0], array[1], array[2]
}

func (x *defaultOIDCClient) HandleSignInCallback(ctx context.Context) {
	session := x.Options.Sessions.Start(ctx)

	state := ctx.FormValue("state")
	if storedState := session.Get(x.Options.Session_State); state != storedState {
		ctx.WriteString("state did not match")
		ctx.StatusCode(http.StatusBadRequest)
		return
	}
	code := ctx.FormValue("code")

	httpCtx := gocontext.Background()
	oauth2Token, err := x.OAuth2Config.Exchange(httpCtx, code)
	if err != nil {
		ctx.Write([]byte("Failed to exchange token: " + err.Error()))
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	userInfo, err := x.OIDCProvider.UserInfo(httpCtx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		ctx.Write([]byte("Failed to get userinfo: " + err.Error()))
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	claims := make(map[string]string, 6)
	userInfo.Claims(&claims)

	session.Set(x.Options.Session_ID, claims["sub"])
	session.Set(x.Options.Session_Username, claims["name"])
	session.Set(x.Options.Session_Roles, claims["role"])
	session.Set(x.Options.Session_Level, claims["level"])
	session.Set(x.Options.Session_Status, claims["status"])
	session.Set(x.Options.Session_Email, claims["email"])

	// 保存令牌
	x.Options.SecureCookie.Set(ctx, x.Options.Cookie_AccesToken, oauth2Token.AccessToken)
	if oauth2Token.RefreshToken != "" {
		x.Options.SecureCookie.Set(ctx, x.Options.Cookie_RefreshToken, oauth2Token.RefreshToken, func(o *http.Cookie) {
			o.Expires = time.Now().Add(336 * time.Hour)
		})
	}

	// Todo: 重定向到登录前页面
	ctx.Redirect("/", http.StatusFound)
}

func (x *defaultOIDCClient) HandleSignOutCallback(ctx context.Context) {
	session := x.Options.Sessions.Start(ctx)

	session.Destroy()
	ctx.RemoveCookie(x.Options.Cookie_AccesToken)
	ctx.RemoveCookie(x.Options.Cookie_RefreshToken)
	ctx.RemoveCookie(x.Options.Cookie_Session)

	// Todo: 去Passport注销
	ctx.Redirect("/", http.StatusFound)
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
	if options.ProviderUrl == "" {
		log.Fatal("OIDCClient.Options.ProviderUrl cannot be empty.")
	}
	if options.RedirectURL == "" {
		log.Fatal("OIDCClient.Options.RedirectURL cannot be empty.")
	}

	if len(options.Scopes) == 0 {
		log.Fatal("OIDCClient.Options.Scopes cannot be empty")
	}
	if options.SecureCookie == nil {
		log.Fatal("OIDCClient.Options.SecureCookie cannot be nil")
	}
	if options.PermissionAuditor == nil {
		log.Fatal("OIDCClient.Options.PermissionAuditor cannot be nil")
	}
	if options.Sessions == nil {
		log.Fatal("OIDCClient.Options.Sessions cannot be nil")
	}

	if options.Cookie_AccesToken == "" {
		options.Cookie_AccesToken = ".ACT"
	}
	if options.Cookie_RefreshToken == "" {
		options.Cookie_RefreshToken = ".RFT"
	}
	if options.Cookie_Session == "" {
		options.Cookie_Session = ".USS"
	}

	if options.Session_ID == "" {
		options.Session_ID = "ID"
	}
	if options.Session_Username == "" {
		options.Session_Username = "Username"
	}
	if options.Session_Email == "" {
		options.Session_Email = "Email"
	}
	if options.Session_Roles == "" {
		options.Session_Roles = "Roles"
	}
	if options.Session_Level == "" {
		options.Session_Level = "Level"
	}
	if options.Session_Status == "" {
		options.Session_Status = "Status"
	}

	if options.Session_State == "" {
		options.Session_State = "State"
	}
}
