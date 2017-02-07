package main

import (
	"fmt"
        "github.com/urfave/cli"
        "github.com/sak0/nvwa/ovs"
)

var joinCommand = cli.Command{
	Name:  "join",
	Usage: "add eth to vswitch",
	ArgsUsage: "<ethX>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "bridge, b", 
			Value: "", 
			Usage: "bridge to join or create",
		},
	},
	Action: joinBridge,
}

func joinBridge(c *cli.Context) {
	device := c.Args().First()
    bridge := c.String("bridge")
    addr, _ := getIfaceAddr(device)
    fmt.Println(bridge, device, addr.IP, addr.String())
    ovsdber, _ := ovs.GetOvsdber()
    //fmt.Printf("%v\n", ovsdber)
    ovsdber.InitCache()
    ovsdber.AddBridge(bridge)
	err := interfaceUp(bridge)
	if err != nil {
		fmt.Printf("Failed: enable iface:%s\n", bridge)
	}
	err = delInterfaceIP(device, addr.String())
    if err != nil {
		fmt.Printf("Failed: del iface [%s] addr %s\n", device, "172.16.74.210/24")
	}

    err = setInterfaceIP(bridge, addr.String())
    if err != nil {
		fmt.Printf("Failed: set iface [%s] addr %s\n", bridge, "172.16.74.210/24")
	}
}
