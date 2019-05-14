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
	node := px.GetNode(os.Args[4])

	fmt.Println(node.ListLXC())
	fmt.Println(node.ListVM())

	lxclist, _ := node.ListLXC()
	for _, lxcid := range lxclist {
		lxc := node.GetLXC(lxcid)
		fmt.Println(lxc.Status())
	}

	vmlist, _ := node.ListVM()
	for _, vmid := range vmlist {
		vm := node.GetVM(vmid)
		fmt.Println(vm.Status())
	}
}
