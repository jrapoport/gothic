package api

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/metering"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
)

// GothicClaims is a struct thats used for JWT claims
type GothicClaims struct {
	jwt.StandardClaims
	Username     string                 `json:"username"`
	Email        string                 `json:"email"`
	Confirmed    bool                   `json:"confirmed"`
	Verified     bool                   `json:"verified"`
	AppMetaData  map[string]interface{} `json:"app_metadata"`
	UserMetaData map[string]interface{} `json:"user_metadata"`
}

// AccessTokenResponse represents an OAuth2 success response
type AccessTokenResponse struct {
	Token        string `json:"access_token"`
	TokenType    string `json:"token_type"` // Bearer
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

const useCookieHeader = "x-use-cookie"
const useSessionCookie = "session"

// Token is the endpoint for OAuth access token requests
func (a *API) Token(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	grantType := r.FormValue("grant_type")

	switch grantType {
	case "password":
		return a.ResourceOwnerPasswordGrant(ctx, w, r)
	case "refresh_token":
		return a.RefreshTokenGrant(ctx, w, r)
	default:
		return oauthError("unsupported_grant_type", "")
	}
}

// ResourceOwnerPasswordGrant implements the password grant type flow
func (a *API) ResourceOwnerPasswordGrant(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	password := r.FormValue("password")
	cookie := r.Header.Get(useCookieHeader)

	aud := a.requestAud(ctx, r)
	config := a.getConfig(ctx)

	if config.Recaptcha.Login {
		if err := a.checkRecaptcha(r, config); err != nil {
			return err
		}
	}

	user, err := models.FindUserByEmailAndAudience(a.db, username, aud)
	if err != nil {
		if models.IsNotFoundError(err) {
			return oauthError("invalid_grant", "No user found with that email, or password invalid.")
		}
		return internalServerError("Name error finding user").WithInternalError(err)
	}

	// NOTE: this is commented out to let users who have not confirmed recover.
	//  this will reduce orphan accounts if the user drops off after signup.
	//if !user.IsConfirmed() {
	//	return oauthError("invalid_grant", "Email not confirmed")
	//}

	if !user.Authenticate(password) {
		return oauthError("invalid_grant", "No user found with that email, or password invalid.")
	}

	var token *AccessTokenResponse
	err = a.db.Transaction(func(tx *storage.Connection) error {
		var terr error
		if terr = models.NewAuditLogEntry(tx, user, models.LoginAction, nil); terr != nil {
			return terr
		}
		if terr = triggerEventHooks(ctx, tx, LoginEvent, user, config); terr != nil {
			return terr
		}

		token, terr = a.issueRefreshToken(ctx, tx, user)
		if terr != nil {
			return terr
		}

		if cookie != "" && config.Cookies.Duration > 0 {
			if terr = a.setCookieToken(config, token.Token, cookie == useSessionCookie, w); terr != nil {
				return internalServerError("Failed to set JWT cookie. %s", terr)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	metering.RecordLogin("password", user.ID)
	return sendJSON(w, http.StatusOK, token)
}

// RefreshTokenGrant implements the refresh_token grant type flow
func (a *API) RefreshTokenGrant(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	config := a.getConfig(ctx)
	tokenStr := r.FormValue("refresh_token")
	cookie := r.Header.Get(useCookieHeader)

	if tokenStr == "" {
		return oauthError("invalid_request", "refresh_token required")
	}

	user, token, err := models.FindUserWithRefreshToken(a.db, tokenStr)
	if err != nil {
		if models.IsNotFoundError(err) {
			return oauthError("invalid_grant", "Invalid Refresh Token")
		}
		return internalServerError(err.Error())
	}

	if token.Revoked {
		a.clearCookieToken(w)
		return oauthError("invalid_grant", "Invalid Refresh Token").WithInternalMessage("Possible abuse attempt: %v", r)
	}

	var tokenString string
	var newToken *models.RefreshToken

	err = a.db.Transaction(func(tx *storage.Connection) error {
		var terr error
		if terr = models.NewAuditLogEntry(tx, user, models.TokenRefreshedAction, nil); terr != nil {
			return terr
		}

		newToken, terr = models.GrantRefreshTokenSwap(tx, user, token)
		if terr != nil {
			return internalServerError(terr.Error())
		}

		tokenString, terr = generateAccessToken(user,
			time.Second*time.Duration(config.JWT.Exp),
			config.JWT.Secret,
			config.JWT.SigningMethod())
		if terr != nil {
			return internalServerError("error generating jwt token").WithInternalError(terr)
		}

		if cookie != "" && config.Cookies.Duration > 0 {
			if terr = a.setCookieToken(config, tokenString, cookie == useSessionCookie, w); terr != nil {
				return internalServerError("Failed to set JWT cookie. %s", terr)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	metering.RecordLogin("token", user.ID)
	return sendJSON(w, http.StatusOK, &AccessTokenResponse{
		Token:        tokenString,
		TokenType:    "bearer",
		ExpiresIn:    config.JWT.Exp,
		RefreshToken: newToken.Token,
	})
}

func generateAccessToken(user *models.User, expiresIn time.Duration, secret string, method jwt.SigningMethod) (string, error) {
	claims := &GothicClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.ID.String(),
			Audience:  jwt.ClaimStrings{user.Aud},
			ExpiresAt: jwt.At(time.Now().Add(expiresIn)),
		},
		Username:     user.Username,
		Email:        user.Email,
		Confirmed:    user.IsConfirmed(),
		Verified:     user.IsVerified(),
		AppMetaData:  user.AppMetaData,
		UserMetaData: user.UserMetaData,
	}

	token := jwt.NewWithClaims(method, claims)
	return token.SignedString([]byte(secret))
}

func (a *API) issueRefreshToken(ctx context.Context, conn *storage.Connection, user *models.User) (*AccessTokenResponse, error) {
	config := a.getConfig(ctx)

	now := time.Now()
	user.LastSignInAt = &now

	var tokenString string
	var refreshToken *models.RefreshToken

	err := conn.Transaction(func(tx *storage.Connection) error {
		var terr error
		refreshToken, terr = models.GrantAuthenticatedUser(tx, user)
		if terr != nil {
			return internalServerError("Name error granting user").WithInternalError(terr)
		}

		tokenString, terr = generateAccessToken(user,
			time.Second*time.Duration(config.JWT.Exp),
			config.JWT.Secret,
			config.JWT.SigningMethod())
		if terr != nil {
			return internalServerError("error generating jwt token").WithInternalError(terr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &AccessTokenResponse{
		Token:        tokenString,
		TokenType:    "bearer",
		ExpiresIn:    config.JWT.Exp,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (a *API) setCookieToken(config *conf.Configuration, tokenString string, session bool, w http.ResponseWriter) error {
	exp := time.Second * time.Duration(config.Cookies.Duration)
	cookie := &http.Cookie{
		Name:     config.Cookies.Key,
		Value:    tokenString,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	if !session {
		cookie.Expires = time.Now().Add(exp)
		cookie.MaxAge = config.Cookies.Duration
	}

	http.SetCookie(w, cookie)
	return nil
}

func (a *API) clearCookieToken(w http.ResponseWriter) {
	config := a.config
	http.SetCookie(w, &http.Cookie{
		Name:     config.Cookies.Key,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour * 10),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	})
}
