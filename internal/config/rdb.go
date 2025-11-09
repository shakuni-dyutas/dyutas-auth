package config

type RDbConfig struct {
	Host string
	Port string
	Db   string
	User string
	Pw   string
}

const (
	confKeyRDbHost = "AUTH_RDB_HOST"
	confKeyRDbPort = "AUTH_RDB_PORT"
	confKeyRDbDb   = "AUTH_RDB_DB"
	confKeyRDbUser = "AUTH_RDB_USER"
	confKeyRDbPw   = "AUTH_RDB_PASSWORD"
)

var loadRDbConfigOf = LoadConfigOf

func loadRDBConfig() (cnf RDbConfig, unconfigureds []string) {
	unconfigureds = []string{}

	host, ok := loadRDbConfigOf(confKeyRDbHost)
	if !ok {
		unconfigureds = append(unconfigureds, confKeyRDbHost)
	}
	port, ok := loadRDbConfigOf(confKeyRDbPort)
	if !ok {
		unconfigureds = append(unconfigureds, confKeyRDbPort)
	}
	db, ok := loadRDbConfigOf(confKeyRDbDb)
	if !ok {
		unconfigureds = append(unconfigureds, confKeyRDbDb)
	}
	user, ok := loadRDbConfigOf(confKeyRDbUser)
	if !ok {
		unconfigureds = append(unconfigureds, confKeyRDbUser)
	}
	pw, ok := loadRDbConfigOf(confKeyRDbPw)
	if !ok {
		unconfigureds = append(unconfigureds, confKeyRDbPw)
	}

	cnf = RDbConfig{
		Host: host,
		Port: port,
		Db:   db,
		User: user,
		Pw:   pw,
	}

	return cnf, unconfigureds
}
