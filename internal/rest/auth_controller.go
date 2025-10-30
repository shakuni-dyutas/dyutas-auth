package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shakuni-dyutas/dyutas-auth/internal/db"
)

func InitAuthControllerWith(db *db.Conn, router gin.IRouter) {
	if db == nil {
		panic("database isn't configured while initializing Auth controller")
	}
	if router == nil {
		panic("gin router isn't configured while initializing Auth controller")
	}

	RegisterHandlersWithOptions(router, &AuthController{db: db}, GinServerOptions{BaseURL: "/auth"})
}

type AuthController struct {
	db *db.Conn
}

type GoogleAPIOAuthTokenRequestDTO struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

type GoogleAPIOAuthTokenResponseDTO struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IdToken      string `json:"id_token"`
}
type GoogleIdTokenPayload struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Iat           int    `json:"iat"`
	Exp           int    `json:"exp"`
}

func (c *GoogleIdTokenPayload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{c.Aud}, nil
}

func (c *GoogleIdTokenPayload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(int64(c.Exp), 0)}, nil
}

func (c *GoogleIdTokenPayload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(int64(c.Iat), 0)}, nil
}

func (c *GoogleIdTokenPayload) GetIssuer() (string, error) {
	return c.Iss, nil
}

func (c *GoogleIdTokenPayload) GetSubject() (string, error) {
	return c.Sub, nil
}

func (c *GoogleIdTokenPayload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (ac *AuthController) GoogleSignIn(ctx *gin.Context) {
	var reqBody GoogleSignInRequestDTO
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to bind request body - " + err.Error()})
		return
	}
	url := "https://oauth2.googleapis.com/token"
	tokenReq := GoogleAPIOAuthTokenRequestDTO{
		Code:         reqBody.Code,
		ClientId:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectUri:  "https://local.dyutas.com:8010",
		GrantType:    "authorization_code",
	}

	jsonReq, err := json.Marshal(tokenReq)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to marshal token request - " + err.Error()})
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to create HTTP request - " + err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to make HTTP request to Google OAuth - " + err.Error()})
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to read response body from Google OAuth - " + err.Error()})
		return
	}

	var respBody GoogleAPIOAuthTokenResponseDTO
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to unmarshal Google OAuth response - " + err.Error()})
		return
	}

	var idTokenJWT = respBody.IdToken

	jwtParser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))

	token, _, err := jwtParser.ParseUnverified(idTokenJWT, &GoogleIdTokenPayload{})
	if err != nil {
		ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to parse Google ID token JWT - " + err.Error()})
		return
	}

	var claims GoogleIdTokenPayload
	if c, ok := token.Claims.(*GoogleIdTokenPayload); ok {
		claims = *c
	}

	existingUser, err := ac.db.Qrs.GetUserByGoogleID(ctx, claims.Sub)
	if err == nil {
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"tid": "nowimintownbreakitdownthinkingofmakinganewsound",
			"ucd": existingUser.Code,
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		})

		auth_jwt_key := os.Getenv("AUTH_APP_JWT_KEY")
		auth_jwt_key_bytes := []byte(auth_jwt_key)

		refTkstr, err := refreshToken.SignedString(auth_jwt_key_bytes)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to sign refresh token - " + err.Error()})
			return
		}

		ctx.SetCookie("rtk", refTkstr, 60*60*24*7, "/", "local-api.dyutas.com", true, true)

		resp := GoogleSignInResponseDTO{
			AccessToken: refTkstr,
			User: UserDTO{
				Email:           existingUser.Email,
				Nickname:        existingUser.Username.String,
				ProfileImageUrl: existingUser.ProfileImageUrl.String,
				UserId:          existingUser.Code,
			},
		}

		ctx.JSON(http.StatusOK, resp)
		return
	}
	if err != pgx.ErrNoRows {
		ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to get user by Google ID - " + err.Error()})
		return
	}

	random10charstringInASCII := []byte{}
	for i := 0; i < 10; i++ {
		random10charstringInASCII = append(random10charstringInASCII, byte(rand.Intn(10)))
	}
	random10charstring := string(random10charstringInASCII)

	user, err := ac.db.Qrs.CreateUser(ctx, db.CreateUserParams{
		Code:            random10charstring,
		GoogleID:        claims.Sub,
		Email:           claims.Email,
		ProfileImageUrl: pgtype.Text{String: claims.Picture, Valid: true},
		Username:        pgtype.Text{String: claims.Name, Valid: true},
		SignedUpAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})

	if err != nil {
		ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to create new user - " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

func (c *AuthController) RefreshAuthentication(ctx *gin.Context) {
	cookie, err := ctx.Cookie("rtk")
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	jwtParser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))
	token, _, err := jwtParser.ParseUnverified(cookie, &jwt.MapClaims{})
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	var claims jwt.MapClaims
	if c, ok := token.Claims.(*jwt.MapClaims); ok {
		claims = *c
	}

	thetokenmapobject := map[string]interface{}{
		"tid": claims["tid"],
		"ucd": claims["ucd"],
		"exp": claims["exp"],
	}

	thetokenmapobject["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(thetokenmapobject))

	auth_jwt_key := os.Getenv("AUTH_APP_JWT_KEY")
	if auth_jwt_key == "" {
		ctx.JSON(400, gin.H{"error": "AUTH_APP_JWT_KEY is not set"})
		return
	}

	auth_jwt_key_bytes := []byte(auth_jwt_key)

	newRefreshTokenString, err := newRefreshToken.SignedString(auth_jwt_key_bytes)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("rtk", newRefreshTokenString, 60*60*24*7, "/", "local-api.dyutas.com", true, true)

	ctx.JSON(http.StatusCreated, gin.H{"accessToken": newRefreshTokenString})
}

func (c *AuthController) SignOut(ctx *gin.Context) {
	ctx.SetCookie("rtk", "", -1, "/", "local-api.dyutas.com", true, true)
	ctx.JSON(http.StatusNoContent, nil)
}
