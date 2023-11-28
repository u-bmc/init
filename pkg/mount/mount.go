// SPDX-License-Identifier: BSD-3-Clause

package mount

import (
	"log"
	"strings"

	"golang.org/x/sys/unix"
)

type Mount struct {
	Source       string
	Destination  string
	Type         string
	GenericFlags uintptr
	Flags        []string
}

type Mounter struct {
	config config
}

func NewMounter(options ...Option) *Mounter {
	cfg := config{}

	for _, opt := range options {
		opt.apply(&cfg)
	}

	return &Mounter{
		config: cfg,
	}
}

func (m *Mounter) MountAll() {
	for _, m := range m.config.mounts {
		if err := unix.Mount(m.Source, m.Destination, m.Type, 0, strings.Join(m.Flags, ",")); err != nil {
			log.Printf("Unable to mount file system %s\n", m.Destination)
		}
	}
}
