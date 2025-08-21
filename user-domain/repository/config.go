package repository

type DBConfig interface {
	Load(filenames ...string) error
	GetDBName() string
	GetHostName() string
	GetPort() string
	GetPassword() string
	GetUser() string
	GetSSLMode() string
}
