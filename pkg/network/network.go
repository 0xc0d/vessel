package network

import (
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"net"
)

func SetupBridge(name string) error {
	// Create Bridge if does not exist
	bridge, err := netlink.LinkByName(name)
	if err != nil {
		linkAttrs := netlink.NewLinkAttrs()
		linkAttrs.Name = name
		bridge = &netlink.Bridge{LinkAttrs: linkAttrs}
		if err := netlink.LinkAdd(bridge); err != nil {
			return err
		}
	}

	// Add IP address if there is not
	addrList, err := netlink.AddrList(bridge, 0)
	if err != nil {
		return err
	}
	if len(addrList) < 1 {
		IP := "172.30.0.1/16"
		addr, err := netlink.ParseAddr(IP)
		if err != nil {
			return err
		}
		if err := netlink.AddrAdd(bridge, addr); err != nil {
			return err
		}
	}

	// Setup the Bridge
	if err := netlink.LinkSetUp(bridge); err != nil {
		return err
	}
	return nil
}

func SetupVirtualEthernet(name, peer string) error {
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = name
	vth := &netlink.Veth{
		LinkAttrs: linkAttrs,
		PeerName:  peer,
	}
	if err := netlink.LinkAdd(vth); err != nil {
		return err
	}
	return netlink.LinkSetUp(vth)
}

func LinkSetMaster(linkName, masterName string) error {
	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return errors.Wrapf(err, "can't find link %s", linkName)
	}
	masterLink, err := netlink.LinkByName(masterName)
	if err != nil {
		return errors.Wrapf(err, "can't find link %s", masterName)
	}
	if err := netlink.LinkSetMaster(link, masterLink); err != nil {
		return err
	}
	return nil
}

func LinkAddGateway(linkName, gatewayIP string) error {
	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return err
	}

	newRoute := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Scope:     netlink.SCOPE_UNIVERSE,
		Gw:        net.ParseIP(gatewayIP),
	}
	return netlink.RouteAdd(newRoute)
}

func LinkAddAddr(linkName, IP string) error {
	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return err
	}
	addr, err := netlink.ParseAddr(IP)
	if err != nil {
		return errors.Wrapf(err, "can't parse %s", IP)
	}
	return netlink.AddrAdd(link, addr)
}

func LinkSetup(linkName string) error {
	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return err
	}
	return netlink.LinkSetUp(link)
}

func LinkRename(old, new string) error {
	link, err := netlink.LinkByName(old)
	if err != nil {
		return err
	}
	return netlink.LinkSetName(link, new)
}


func IPExists(ip net.IP) (bool, error) {
	linkList, err := netlink.AddrList(nil, 0)
	if err != nil {
		return false, err
	}
	for _, link := range linkList {
		if link.IP.Equal(ip) {
			return true, nil
		}
	}
	return false, nil
}
