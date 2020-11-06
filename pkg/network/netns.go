package network

import (
	"github.com/0xc0d/vessel/pkg/filesystem"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

type Unsetter func() error

func MountNewNetworkNamespace(nsTarget string) (filesystem.Unmounter, error) {
	_, err := os.OpenFile(nsTarget, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_EXCL, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create target file")
	}

	// store current network namespace
	file, err := os.OpenFile("/proc/self/ns/net", os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := syscall.Unshare(syscall.CLONE_NEWNET); err != nil {
		return nil, errors.Wrap(err, "unshare syscall failed")
	}
	mountPoint := filesystem.MountOption{
		Source: "/proc/self/ns/net",
		Target: nsTarget,
		Type:   "bind",
		Flag:   syscall.MS_BIND,
	}
	unmount, err := filesystem.Mount(mountPoint)
	if err != nil {
		return unmount, err
	}

	// reset previous network namespace
	if err := unix.Setns(int(file.Fd()), syscall.CLONE_NEWNET); err != nil {
		return unmount, errors.Wrap(err, "setns syscall failed: ")
	}

	return unmount, nil
}

func SetNetNSByFile(filename string) (Unsetter, error) {
	currentNS, err := os.OpenFile("/proc/self/ns/net", os.O_RDONLY, 0)
	unsetFunc := func() error {
		defer currentNS.Close()
		if err != nil {
			return err
		}
		return unix.Setns(int(currentNS.Fd()), syscall.CLONE_NEWNET)
	}

	netnsFile, err := os.OpenFile(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return unsetFunc, errors.Wrap(err, "unable to open network namespace file")
	}
	defer netnsFile.Close()
	if err := unix.Setns(int(netnsFile.Fd()), syscall.CLONE_NEWNET); err != nil {
		return unsetFunc, errors.Wrap(err, "unset syscall failed")
	}
	return unsetFunc, err
}

func LinkSetNsByFile(filename, linkName string) error {
	netnsFile, err := os.OpenFile(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return errors.Wrap(err, "unable to open netns file")
	}
	defer netnsFile.Close()

	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return err
	}
	return netlink.LinkSetNsFd(link, int(netnsFile.Fd()))
}
