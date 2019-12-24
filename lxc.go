package goproxmox

type LXC struct {
	OSTemplate string `json:"ostemplate"`
	VMID       string `json:"vmid"`

	Arch               string  `json:"arch,omitempty"`
	BWLimit            float64 `json:"bwlimit,omitempty"`
	CMode              string  `json:"cmode,omitempty"`
	Console            int     `json:"console,omitempty"`
	Cores              uint    `json:"cores,omitempty"`
	CPULimit           float64 `json:"cpulimit,omitempty"`
	CPUUnits           int64   `json:"cpuunits,omitempty"`
	Description        string  `json:"description,omitempty"`
	Features           string  `json:"features,omitempty"`
	Force              int     `json:"force,omitempty"`
	Hookscript         string  `json:"hookscript,omitempty"`
	Hostname           string  `json:"hostname,omitempty"`
	IgnoreUnpackErrors int     `json:"ignore_unpack_errors,omitempty" override:"ignore-unpack-errors"`
	Lock               string  `json:"lock,omitempty"`
	Memory             uint    `json:"memory,omitempty"`

	MP []string `json:"mp"`

	Nameserver   string `json:"nameserver,omitempty"`
	SearchDomain string `json:"searchdomain,omitempty"`

	Net []VBaseNICSettings `json:"net"`

	OnBoot     int    `json:"onboot,omitempty"`
	OsType     string `json:"ostype,omitempty"`
	Password   string `json:"password,omitempty"`
	Pool       string `json:"pool,omitempty"`
	Protection int    `json:"protection,omitempty"`
	Restore    int    `json:"restore,omitempty"`
	RootFS     string `json:"rootfs,omitempty"`

	SSH_Public_Keys string `json:"ssh_public_keys,omitempty" override:"ssh-public-keys"`
	Start           int    `json:"start,omitempty"`
	Startup         string `json:"startup,omitempty"`
	Storage         string `json:"storage,omitempty"`
	Swap            uint   `json:"swap,omitempty"`
	Template        int    `json:"template,omitempty"`
	TTY             int    `json:"tty,omitempty"`
	Unique          int    `json:"unique,omitempty"`
	Unprivileged    int    `json:"unprivileged"`
}
