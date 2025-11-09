package googleadpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type GoogleAdptConfig struct {
	ClientId     string
	ClientSecret string
}

func New(config GoogleAdptConfig) GoogleAdpt {
	return &GoogleAdptImpl{
		clientId:     config.ClientId,
		clientSecret: config.ClientSecret,
	}
}

type GoogleAdptImpl struct {
	clientId     string
	clientSecret string
}

const googleSignURL = "https://oauth2.googleapis.com/token"
const gooleSignGrantType = "authorization_code"

// currently meaningless
const googleSignRedirectURI = "https://local.dyutas.com:8010"

func (adpt *GoogleAdptImpl) VerifyAuthCode(ctx context.Context, googleAuthCode string) (*GoogleUserInfo, error) {
	googleIdToken, err := adpt.requestCodeAuthentication(googleAuthCode)
	if err != nil {
		return nil, fmt.Errorf("failed to request code authentication: %w", err)
	}

	googleUserInfo := GoogleUserInfo{
		Sub:             googleIdToken.Sub,
		Email:           googleIdToken.Email,
		ProfileImageURL: googleIdToken.Picture,
	}

	return &googleUserInfo, nil
}

func (adpt *GoogleAdptImpl) requestCodeAuthentication(googleAuthCode string) (*GoogleIdTokenPayload, error) {
	signHTTPReq, err := adpt.generateGoogleSignHTTPRequest(googleAuthCode)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTTP request body for Google sign: %w", err)
	}

	signHTTPRespBodyDTO, err := adpt.requestToGoogleAuthAPIWith(signHTTPReq)
	if err != nil {
		return nil, fmt.Errorf("failed to request to Google Auth API: %w", err)
	}

	googleIdToken, err := adpt.parseGoogleIdTokenJWT(signHTTPRespBodyDTO.IdToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Google ID token JWT: %w", err)
	}

	return googleIdToken, nil
}

func (adpt *GoogleAdptImpl) generateGoogleSignHTTPRequest(googleAuthCode string) (*http.Request, error) {
	signReq := GoogleSignAPIRequestDTO{
		Code:         googleAuthCode,
		ClientId:     adpt.clientId,
		ClientSecret: adpt.clientSecret,
		RedirectUri:  googleSignRedirectURI,
		GrantType:    gooleSignGrantType,
	}

	signReqBytes, err := json.Marshal(signReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Google sign request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, googleSignURL, bytes.NewBuffer(signReqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request for Google sign: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	return httpReq, nil
}

func (adpt *GoogleAdptImpl) requestToGoogleAuthAPIWith(reqBody *http.Request) (*GoogleSignAPIResponseDTO, error) {
	signHTTPResp, err := http.DefaultClient.Do(reqBody)
	if err != nil {
		return nil, err
	}
	defer signHTTPResp.Body.Close()

	signHTTPRespBody, err := io.ReadAll(signHTTPResp.Body)
	if err != nil {
		return nil, err
	}

	var signHTTPRespBodyDTO GoogleSignAPIResponseDTO
	err = json.Unmarshal(signHTTPRespBody, &signHTTPRespBodyDTO)
	if err != nil {
		return nil, err
	}

	return &signHTTPRespBodyDTO, nil
}

// TODO id token verification logic needed before open service
func (adpt *GoogleAdptImpl) parseGoogleIdTokenJWT(idTokenJWT string) (*GoogleIdTokenPayload, error) {
	jwtParser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))

	token, _, err := jwtParser.ParseUnverified(idTokenJWT, &GoogleIdTokenPayload{})
	if err != nil {
		return nil, err
	}

	var claims GoogleIdTokenPayload
	if c, ok := token.Claims.(*GoogleIdTokenPayload); ok {
		claims = *c
	}

	return &claims, nil
}
