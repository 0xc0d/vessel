package container

import (
	"fmt"
	"github.com/0xc0d/vessel/pkg/filesystem"
	"github.com/0xc0d/vessel/pkg/network"
	"path/filepath"
)

const netnsPath = "/var/run/vessel/netns"

func (c *Container) SetupNetwork(bridge string) (filesystem.Unmounter, error) {
	nsMountTarget := filepath.Join(netnsPath, c.Digest)
	vethName := fmt.Sprintf("veth%.7s", c.Digest)

	if err := network.SetupVirtualEthernet(vethName, "eth0z", bridge); err != nil {
		return nil, err
	}
	unmount, err := network.MountNewNetworkNamespace(nsMountTarget)
	if err != nil {
		return unmount, err
	}
	if err := network.LinkSetNsFile(nsMountTarget, "eth0z"); err != nil {
		return unmount, err
	}
	unset, err := network.SetNetNSByFile(nsMountTarget)
	if err != nil {
		return unmount, nil
	}
	defer unset()
	if err := network.LinkAddAddr("eth0z", "172.30.0.2/16"); err != nil {
		return unmount, err
	}
	if err := network.LinkSetup("eth0z"); err != nil {
		return unmount, err
	}
	if err := network.LinkAddGateway("eth0z", "172.30.0.1"); err != nil {
		return unmount, err
	}
	if err := network.LinkSetup("lo"); err != nil {
		return unmount, err
	}

	return unmount, nil
}

func (c *Container) SetNetworkNamespace() (network.Unsetter, error) {
	netns := filepath.Join(netnsPath, c.Digest)
	return network.SetNetNSByFile(netns)
}
