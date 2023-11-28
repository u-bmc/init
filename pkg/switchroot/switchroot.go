// SPDX-License-Identifier: BSD-3-Clause

package switchroot

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func SwitchRoot(newroot string, init string) {
	log.Println("Moving mounts")

	mounts := []string{"/dev", "/proc", "/sys", "/run"}
	for _, mount := range mounts {
		moveMount(mount, filepath.Join(newroot, mount))
	}

	oldroot, err := os.Open("/")
	if err != nil {
		log.Printf("switch_root: failed to open /: %v", err)
	}
	defer oldroot.Close()

	moveMount(newroot, "/")

	if err := unix.Chroot("."); err != nil {
		log.Printf("switch_root: failed to call chroot: %v", err)
	}

	recursiveDelete(int(oldroot.Fd()))

	log.Printf("Executing %s", init)
	if err := unix.Exec(init, []string{init}, []string{}); err != nil {
		log.Printf("switch_root: failed to exec %s: %v", init, err)
	}
}

func moveMount(oldPath, newPath string) {
	if err := unix.Mount(oldPath, newPath, "", unix.MS_MOVE, ""); err != nil {
		log.Printf("switch_root: failed to move mount %s: %v", oldPath, err)
	}
}

func recursiveDelete(fd int) {
	parentDev, err := getDev(fd)
	if err != nil {
		log.Printf("switch_root: unable to get underlying dev for dir: %v", err)
		return
	}

	dir := os.NewFile(uintptr(fd), "__ignored__")
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		log.Printf("switch_root: unable to read dir %s: %v", dir.Name(), err)
		return
	}

	for _, name := range names {
		recusiveDeleteInner(fd, parentDev, name)
	}
}

func recusiveDeleteInner(parentFd int, parentDev uint64, childName string) {
	childFd, err := unix.Openat(parentFd, childName, unix.O_DIRECTORY|unix.O_NOFOLLOW, unix.O_RDWR)
	if err != nil {
		if err := unix.Unlinkat(parentFd, childName, 0); err != nil {
			log.Printf("switch_root: unable to remove file %s: %v", childName, err)
		}
	} else {
		defer unix.Close(childFd)

		if childFdDev, err := getDev(childFd); err != nil {
			log.Printf("switch_root: unable to get underlying dev for dir: %s: %v", childName, err)
			return
		} else if childFdDev != parentDev {
			return
		}

		recursiveDelete(childFd)

		if err := unix.Unlinkat(parentFd, childName, unix.AT_REMOVEDIR); err != nil {
			log.Printf("switch_root: unable to remove dir %s: %v", childName, err)
		}
	}
}

func getDev(fd int) (uint64, error) {
	var stat unix.Stat_t
	if err := unix.Fstat(fd, &stat); err != nil {
		return 0, err
	}
	return uint64(stat.Dev), nil
}
