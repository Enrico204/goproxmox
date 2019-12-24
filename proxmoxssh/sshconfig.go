package proxmoxssh

type Config struct {
	Hostname      string
	Port          int
	User          string
	PrivateKey    []byte
	Password      string
	HostPublicKey []byte
}
