// SPDX-License-Identifier: BSD-3-Clause

package keyring

import (
	"log"

	"golang.org/x/sys/unix"
)

func AddUbifsAuthKey(key []byte) {
	id, err := unix.AddKey("logon", "ubifs:auth", key, unix.KEY_SPEC_SESSION_KEYRING)
	if err != nil {
		log.Printf("Unable to add ubifs auth key: %v", err)
	}

	log.Printf("Added ubifs auth key with id %d", id)
}
