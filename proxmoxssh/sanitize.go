package proxmoxssh

import "strings"

func Sanitize(cmd string) string {
	cmd = strings.ReplaceAll(cmd, "\"", "")
	cmd = strings.ReplaceAll(cmd, "$", "")
	cmd = strings.ReplaceAll(cmd, "`", "")
	cmd = strings.ReplaceAll(cmd, "!", "")
	cmd = strings.ReplaceAll(cmd, "\n", "")
	cmd = strings.ReplaceAll(cmd, "\r", "")
	cmd = strings.ReplaceAll(cmd, "\\", "")
	return cmd
}
