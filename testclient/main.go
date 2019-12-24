package main

import (
	"fmt"
	"gitlab.com/Enrico204/goproxmox"
	"net"
	"os"
	"time"
)

func main() {
	px, err := goproxmox.NewClient(os.Args[1], false, "")
	if err != nil {
		panic(err)
	}
	err = px.Login(os.Args[2], os.Args[3])
	if err != nil {
		panic(err)
	}

	node := px.GetNode(os.Args[4])

	fmt.Println(node.GetVM("100").GuestPing())
	os.Exit(0)

	_, err = node.NewLXC(goproxmox.LXC{
		OSTemplate:   "local:vztmpl/debian-10.0-standard_10.0-1_amd64.tar.gz",
		Password:     "Passw0rd.1",
		Hostname:     "Hostname",
		Unprivileged: 1,
		RootFS:       "local-lvm:8",
		Cores:        1,
		Memory:       512,
		Swap:         512,
		Net: []goproxmox.VBaseNICSettings{
			{
				Id:       0,
				Bridge:   "vmbr0",
				Tag:      1234,
				Firewall: false,
				Name:     "eth0",
				DHCPv4:   true,
			},
			{
				Id:       1,
				Bridge:   "vmbr0",
				Tag:      1235,
				Firewall: false,
				Name:     "eth1",
				Manualv4: true,
			},
			{
				Id:       2,
				Bridge:   "vmbr0",
				Tag:      1236,
				Firewall: false,
				Name:     "eth2",
				IPv4: net.IPNet{
					IP:   net.ParseIP("1.2.3.4"),
					Mask: net.IPv4Mask(255, 255, 255, 0),
				},
				Gateway4: net.ParseIP("1.1.1.1"),
			},
		},
	}, 5*time.Minute)
	fmt.Println(err)
}
