package proxmoxssh

import "fmt"

func PctSetUserPassword(cfg Config, containerid int, username string, password string) error {
	_, err := SimpleRemoteRun(cfg, fmt.Sprintf("pct exec %d chpasswd", containerid), username+":"+password)
	return err
}
