package main

import (
	"fmt"
	"gitlab.com/Enrico204/goproxmox"
	"os"
)

func main() {
	px := goproxmox.NewClient(os.Args[1], false)
	err := px.Login(os.Args[2], os.Args[3])
	if err != nil {
		panic(err)
	}
	//node := px.GetNode(os.Args[4])

	//networks, err := node.ListNetworks()
	//fmt.Println(err)
	//for _, net := range networks {
	//	fmt.Print(net.IFace, " (", net.Type, ")")
	//	if net.Address != nil {
	//		fmt.Print(" -> ", *net.Address)
	//	}
	//	if net.Netmask != nil {
	//		fmt.Print("/", *net.Netmask)
	//	}
	//	if net.Gateway != nil {
	//		fmt.Print(" , ", *net.Gateway)
	//	}
	//
	//	if net.OVS_Bridge != nil {
	//		fmt.Print(" bridged to ", *net.OVS_Bridge)
	//	}
	//	if net.OVS_Ports != nil {
	//		fmt.Print(" bridging ports: ", *net.OVS_Ports)
	//	}
	//	fmt.Println()
	//}

	//comments := "pippo"
	//err = node.CreateNetworkConfig(goproxmox.Network{
	//	IFace: "vmbr10",
	//	Type: "bridge",
	//	Comments: &comments,
	//})
	//err = node.ReloadNetworkConfig()
	//err = node.RevertNetworkChanges()
	fmt.Println(err)
}
