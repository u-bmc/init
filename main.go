// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"flag"
	"log"

	"github.com/u-bmc/init/pkg/keyring"
	"github.com/u-bmc/init/pkg/mount"
	"github.com/u-bmc/init/pkg/switchroot"
	"golang.org/x/sys/unix"
)

const logo = `
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⣤⣤⣤⣤⣀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⣴⣿⠟⠉⠁⠀⠈⠉⠛⢿⣦⡀
⠀⠀⠀⠀⠀⠀⠀⣠⣿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⣿⣆
⠀⠀⠀⠀⠀⠀⢠⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⡆
⠀⠀⠀⠀⠀⠀⣿⡇⠀⠀⣶⣶⣦⠀⠀⠀⣠⣾⣶⡀⠀⠸⣿
⠀⠀⠀⠀⠀⠀⣿⠀⠀⠘⠛⠀⠟⠀⠀⠀⠻⠁⠙⠃⠀⠀⣿
⠀⠀⠀⠀⠀⠀⣿⡇⠀⠀⠀⠀⣴⡆⠀⢠⣦⠀⠀⠀⠀⢠⣿
⠀⠀⠀⠀⠀⠀⠘⣿⠀⠀⠀⠀⠘⣿⣶⣿⠃⠀⠀⠀⠀⣿⠇
⠀⠀⠀⠀⠀⠀⠀⠙⣿⣾⠿⠀⠀⠀⠀⠀⠀⠀⠿⣿⣾⠟
⠀⠀⠀⠀⠀⠀⢀⣿⠛⠀⠀⠀⢀⣤⣤⣤⡀⠀⠀⠀⠙⣿⡄
⠀⠀⠀⠀⠀⢀⣿⠃⠀⠀⠀⣴⡿⠋⠉⠉⢿⣷⠀⠀⠀⠈⣿⡄
⠀⠀⠀⠀⠀⣼⡏⠀⠀⠀⢠⣿⠀⠀⠀⠀⠀⣿⡇⠀⠀⠀⢸⣿
⠀⠀⠀⠀⠀⣿⡇⠀⠀⠀⢸⣿⠀⠀⠀⠀⠀⣿⡇⠀⠀⠀⢀⣿
⠀⠀⠀⢰⣿⠛⠻⣿⠀⠀⢸⣿⠀⠀⠀⠀⠀⣿⡇⠀⠀⣾⠟⠛⢿⣦
⠀⠀⠀⢿⣇⠀⠀⣿⠇⠀⢸⣿⠀⠀⠀⠀⠀⣿⡇⠀⠀⣿⠀⠀⢠⣿
⠀⠀⠀⠈⠻⣿⡿⠋⠀⠀⢸⣿⠀⠀⠀⠀⠀⣿⡇⠀⠀⠙⠿⣿⠿⠁
⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣼⣿⣤⠀⠀⠀⣠⣿⣧⣄
⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⠁⠀⢻⣷⠀⣸⡟⠀⠈⣿⡆
⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⣄⣀⣾⠏⠀⠸⣷⣄⣠⣿⠃
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠛⠁⠀⠀⠀⠈⠛⠛

⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣶⡄
⠀⣠⣄⠀⢠⣦⠀⠀⠀⠀⣿⣁⣶⣤⡀⠀⢀⣤⣦⣀⣴⣦⡀⠀⠀⣤⣶⣤⡀
⠀⣿⡷⠀⢸⣿⢀⣤⣤⠀⣿⡟⠉⢻⣿⠀⣿⠏⠉⣿⠋⠙⣿⠀⣿⡟⠉⠛⠃
⠀⢻⣷⣀⣾⡟⠀⠉⠉⠀⣿⣧⣀⣼⣿⠀⣿⠀⠀⣿⠀⠀⣿⠀⣿⣧⣀⣴⡆
⠀⠀⠉⠛⠋⠀⠀⠀⠀⠀⠛⠉⠛⠋⠀⠀⠛⠀⠀⠛⠀⠀⠛⠀⠀⠙⠛⠋
`

func main() {
	log.Print(logo)

	key := flag.String("key", "", "Filesystem authentication key")
	rootfs := flag.String("rootfs", "", "Root filesystem")
	flag.Parse()

	log.Println("Mounting all file systems")
	m := mount.NewMounter(
		mount.WithDefaultMounts(),
		mount.WithMount(mount.Mount{
			Source:       *rootfs,
			Destination:  "/newroot",
			Type:         "ubifs",
			GenericFlags: unix.MS_NOSUID | unix.MS_RELATIME | unix.MS_LAZYTIME | unix.MS_RDONLY,
			Flags:        []string{"ro", "auth_key=ubifs:auth", "auth_hash_name=sha256", "chk_data_crc", "bulk_read"},
		}),
	)
	m.MountAll()

	log.Println("Populating kernel keyring")
	keyring.AddUbifsAuthKey([]byte(*key))

	log.Println("Changing into new rootfs")
	switchroot.SwitchRoot("/newroot", "/sbin/operator")
}
