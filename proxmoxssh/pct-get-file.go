package proxmoxssh

import "fmt"

func PctGetFile(cfg Config, containerid int, fname string) (string, error) {
	return SimpleRemoteRun(cfg, fmt.Sprintf("pct exec %d cat %s", containerid, Sanitize(fname)), "")
}
