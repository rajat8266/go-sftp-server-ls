package config

// Config represents information go-sftp-server needs to operate
type Config struct {
	Accounts               []Accounts
	BasePath               string
	SSHKeyPath             string
	Address                string
	Port                   string
	SftpAuthorizedKeysFile string
}

// Account holds specific information for each account we support
type Accounts struct {
	Username string
	Password string
}
