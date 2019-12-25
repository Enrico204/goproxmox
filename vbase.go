package goproxmox

import (
	"time"
)

type VBase interface {
	Id() string
	Type() string

	Status() (*MemberStatus, error)
	Start(timeout time.Duration) error
	Stop(timeout time.Duration) error
	Shutdown(timeout time.Duration) error
	Delete(purge bool, timeout time.Duration) error
	Clone(newhostname string, pool string, full bool, newNodeName string, timeout time.Duration) (string, error)
	//Info() error

	SetNIC(settings VBaseNICSettings) error
	DeleteNIC(id int) error

	WaitForGuest(timeout time.Duration) (bool, error)
	GuestPing() (bool, error)
	GuestExecAsync(cmd string) (uint, error)
	GuestExecStatus(pid uint) (GuestExecResult, error)
	GuestExecSync(cmd string) (GuestExecResult, error)
	GuestSetUserPassword(username string, password string) error
	GuestFileRead(fname string) (string, error)
	GuestFileWrite(fname string, content string) error
}

type vbaseimpl struct {
	vmtype string // Can be "lxc" or "qemu"
	id     string
	node   *nodeImpl
}

func (v *vbaseimpl) Id() string {
	return v.id
}

func (v *vbaseimpl) Type() string {
	return v.vmtype
}
