package config

type AuthConfig struct {
	AppJwtKey          string
	GoogleClientId     string
	GoogleClientSecret string
}

const confKeyAppJWTKey = "AUTH_APP_JWT_KEY"
const confKeyGoogleClientId = "GOOGLE_CLIENT_ID"
const confKeyGoogleClientSecret = "GOOGLE_CLIENT_SECRET"

var loadAuthConfigOf = LoadConfigOf

func loadAuthConfig() (cnf AuthConfig, unconfigureds []string) {
	unconfigureds = []string{}

	appJwtKey, ok := loadAuthConfigOf(confKeyAppJWTKey)
	if !ok {
		unconfigureds = append(unconfigureds, confKeyAppJWTKey)
	}
	googleClientId, ok := loadAuthConfigOf(confKeyGoogleClientId)
	if !ok {
		unconfigureds = append(unconfigureds, confKeyGoogleClientId)
	}
	googleClientSecret, ok := loadAuthConfigOf(confKeyGoogleClientSecret)
	if !ok {
		unconfigureds = append(unconfigureds, confKeyGoogleClientSecret)
	}

	cnf = AuthConfig{
		AppJwtKey:          appJwtKey,
		GoogleClientId:     googleClientId,
		GoogleClientSecret: googleClientSecret,
	}

	return cnf, unconfigureds
}
