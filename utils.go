package main

import (
	"fmt"
	"net"
    "time"
    "os/exec"
    "errors"
	"github.com/vishvananda/netlink"
    log "github.com/Sirupsen/logrus"
)

var (
	cmdTimeout time.Duration = time.Minute
)

func getIfaceAddr(name string) (*net.IPNet, error) {
	iface, err := netlink.LinkByName(name)
	if err != nil {
                fmt.Printf("error")
		return nil, err
	}
	addrs, err := netlink.AddrList(iface, netlink.FAMILY_V4)
	if err != nil {
		return nil, err
	}
	fmt.Printf("addrs: %v\n", addrs)
	if len(addrs) == 0 {
		return nil, fmt.Errorf("Interface %s has no IP addresses", name)
	}
	if len(addrs) > 1 {
		fmt.Printf("Interface [ %v ] has more than 1 IPv4 address. Defaulting to using [ %v ]\n", name, addrs[0].IP)
	}
	return addrs[0].IPNet, nil
}

// Set the IP addr of a netlink interface
func setInterfaceIP(name string, rawIP string) error {
	retries := 2
	var iface netlink.Link
	var err error
	for i := 0; i < retries; i++ {
		iface, err = netlink.LinkByName(name)
		if err == nil {
			break
		}
		log.Debugf("error retrieving new OVS bridge netlink link [ %s ]... retrying", name)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Abandoning retrieving the new OVS bridge link from netlink, Run [ ip link ] to troubleshoot the error: %s", err)
		return err
	}
	ipNet, err := netlink.ParseIPNet(rawIP)
	if err != nil {
		fmt.Printf("setInterfaceIP failed ParseIPNet.")
		return err
	}
	
	addr := &netlink.Addr{IPNet: ipNet, Label: ""}
	return netlink.AddrAdd(iface, addr)
}

func delInterfaceIP(name string, rawIP string) error {
	retries := 2
	var iface netlink.Link
	var err error
	for i := 0; i < retries; i++ {
		iface, err = netlink.LinkByName(name)
		if err == nil {
			break
		}
		log.Debugf("error retrieving new OVS bridge netlink link [ %s ]... retrying", name)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Abandoning retrieving the new OVS bridge link from netlink, Run [ ip link ] to troubleshoot the error: %s", err)
		return err
	}
	ipNet, err := netlink.ParseIPNet(rawIP)
	if err != nil {
		return err
	}
	addr := &netlink.Addr{IPNet: ipNet, Label: ""}
	return netlink.AddrDel(iface, addr)
}

func Execute(binary string, args []string) (string, error) {
	var output []byte
	var err error
	cmd := exec.Command(binary, args...)
	fmt.Println(cmd)
	done := make(chan struct{})

	go func() {
		output, err = cmd.CombinedOutput()
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(cmdTimeout):
		if cmd.Process != nil {
			fmt.Printf("Kill")
			cmd.Process.Kill()
		}
		return "", fmt.Errorf("Timeout executing: %v %v, output %v, error %v", binary, args, string(output), err)
	}

	if err != nil {
		return "", fmt.Errorf("Failed to execute: %v %v, output %v, error %v", binary, args, string(output), err)
	}
	return string(output), nil
}

func interfaceUp(name string) error {
    for i := 0; i < 5; i++ {
		iface, err := netlink.LinkByName(name)
		if err == nil {
			return netlink.LinkSetUp(iface)
		}
        fmt.Printf("can't find link [%s] , retry... \n", name)
        time.Sleep(time.Millisecond * 100)
	}
    return errors.New("Error retrieving a link")
}
