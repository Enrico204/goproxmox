package proxmoxssh

import "fmt"

func PctPutFile(cfg Config, containerid int, fname string, content string) error {
	_, err := SimpleRemoteRun(cfg, fmt.Sprintf("pct exec %d tee %s", containerid, Sanitize(fname)), content)
	return err
}
