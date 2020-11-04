package network

import (
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"net"
)

func SetupBridge(name, IP string) error {
	bridge, err := netlink.LinkByName(name)
	if err != nil {
		linkAttrs := netlink.NewLinkAttrs()
		linkAttrs.Name = name
		bridge = &netlink.Bridge{LinkAttrs: linkAttrs}
		if err := netlink.LinkAdd(bridge); err != nil {
			return err
		}
		addr, err := netlink.ParseAddr(IP)
		if err != nil {
			return err
		}
		if err := netlink.AddrAdd(bridge, addr); err != nil {
			return err
		}
	}
	if err := netlink.LinkSetUp(bridge); err != nil {
		return err
	}
	return nil
}

func SetupVirtualEthernet(name, peer, master string) error {
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = name
	vth := &netlink.Veth{
		LinkAttrs: linkAttrs,
		PeerName:  peer,
	}
	if err := netlink.LinkAdd(vth); err != nil {
		return err
	}
	if err := netlink.LinkSetUp(vth); err != nil {
		return err
	}
	masterLink, err := netlink.LinkByName(master)
	if err != nil {
		return errors.Wrapf(err, "can't find link %s", master)
	}
	if err := netlink.LinkSetMaster(vth, masterLink); err != nil {
		return err
	}
	return nil
}

func LinkAddGateway(linkName, gateway string) error {
	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return err
	}

	newRoute := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Scope:     netlink.SCOPE_UNIVERSE,
		Gw:        net.ParseIP(gateway),
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
