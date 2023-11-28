// SPDX-License-Identifier: BSD-3-Clause

package mount

import "golang.org/x/sys/unix"

type config struct {
	mounts []Mount
}

type Option interface {
	apply(*config)
}

type mountOption struct {
	mount Mount
}

func (o *mountOption) apply(c *config) {
	c.mounts = append(c.mounts, o.mount)
}

func WithMount(m Mount) Option {
	return &mountOption{
		mount: m,
	}
}

type mountsOption struct {
	mounts []Mount
}

func (o *mountsOption) apply(c *config) {
	c.mounts = append(c.mounts, o.mounts...)
}

func WithMounts(m []Mount) Option {
	return &mountsOption{
		mounts: m,
	}
}

const (
	DefaultFlagMask uintptr = unix.MS_NOEXEC | unix.MS_NOSUID | unix.MS_NODEV | unix.MS_RELATIME | unix.MS_LAZYTIME
	DefaultFlags            = "rw,noexec,nosuid,nodev,relatime,lazytime"
)

func WithDefaultMounts() Option {
	return &mountsOption{
		mounts: []Mount{
			{
				Source:       "proc",
				Destination:  "/proc",
				Type:         "proc",
				GenericFlags: DefaultFlagMask,
				Flags:        []string{DefaultFlags},
			},
			{
				Source:       "sys",
				Destination:  "/sys",
				Type:         "sysfs",
				GenericFlags: DefaultFlagMask,
				Flags:        []string{DefaultFlags},
			},
			{
				Source:       "securityfs",
				Destination:  "/sys/kernel/security",
				Type:         "securityfs",
				GenericFlags: DefaultFlagMask,
				Flags:        []string{DefaultFlags},
			},
			{
				Source:       "cgroup2",
				Destination:  "/sys/fs/cgroup",
				Type:         "cgroup2",
				GenericFlags: DefaultFlagMask,
				Flags:        []string{DefaultFlags, "nsdelegate", "memory_recursiveprot"},
			},
			{
				Source:       "bpf",
				Destination:  "/sys/fs/bpf",
				Type:         "bpf",
				GenericFlags: DefaultFlagMask,
				Flags:        []string{DefaultFlags, "mode=700"},
			},
			{
				Source:       "configfs",
				Destination:  "/sys/kernel/config",
				Type:         "configfs",
				GenericFlags: DefaultFlagMask,
				Flags:        []string{DefaultFlags},
			},
			{
				Source:       "debugfs",
				Destination:  "/sys/kernel/debug",
				Type:         "debugfs",
				GenericFlags: DefaultFlagMask,
				Flags:        []string{DefaultFlags, "mode=700"},
			},
			{
				Source:       "tracefs",
				Destination:  "/sys/kernel/tracing",
				Type:         "tracefs",
				GenericFlags: DefaultFlagMask,
				Flags:        []string{DefaultFlags, "mode=700"},
			},
			{
				Source:       "dev",
				Destination:  "/dev",
				Type:         "devtmpfs",
				GenericFlags: unix.MS_NOSUID | unix.MS_RELATIME | unix.MS_LAZYTIME,
				Flags:        []string{"rw", "nosuid", "relatime", "lazytime", "mode=755"},
			},
			{
				Source:       "shm",
				Destination:  "/dev/shm",
				Type:         "tmpfs",
				GenericFlags: unix.MS_NOSUID | unix.MS_NODEV | unix.MS_RELATIME | unix.MS_LAZYTIME,
				Flags:        []string{"rw", "nosuid", "nodev", "relatime", "lazytime"},
			},
			{
				Source:       "devpts",
				Destination:  "/dev/pts",
				Type:         "devpts",
				GenericFlags: unix.MS_NOSUID | unix.MS_RELATIME | unix.MS_LAZYTIME,
				Flags:        []string{"rw", "nosuid", "relatime", "lazytime", "mode=620", "gid=5", "ptmxmode=000"},
			},
			{
				Source:       "run",
				Destination:  "/run",
				Type:         "tmpfs",
				GenericFlags: unix.MS_NOSUID | unix.MS_NODEV | unix.MS_RELATIME | unix.MS_LAZYTIME,
				Flags:        []string{"rw", "nosuid", "nodev", "relatime", "lazytime", "mode=755"},
			},
		},
	}
}
