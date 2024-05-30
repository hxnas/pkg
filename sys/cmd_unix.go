//go:build !darwin && !windows
// +build !darwin,!windows

package sys

import (
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func SetCredential(c *exec.Cmd, uid, gid uint32) {
	if uid != 0 || gid != 0 {
		attrInit(c).Credential = &syscall.Credential{Uid: uid, Gid: gid, NoSetGroups: true}
	}
}

func SetNewNS(c *exec.Cmd) {
	attrInit(c).Cloneflags |= syscall.CLONE_NEWNS
}

func SetPdeathsig(c *exec.Cmd, sig syscall.Signal) {
	if sig != 0 {
		attrInit(c).Pdeathsig = sig
	}
}

func SetChroot(c *exec.Cmd, root string) {
	if root != "" && root != "/" {
		attrInit(c).Chroot = root
	}
}

func UserLookup(userOrId, groupOrId string) (uid, gid uint32, err error) {
	var u *user.User
	var g *user.Group

	if userOrId != "" {
		u, err = user.Lookup(userOrId)

		if err != nil {
			u, err = user.LookupId(userOrId)
		}

		if err != nil {
			uid, err = parseUint(userOrId)
		}

		if err != nil {
			return
		}
	}

	if groupOrId != "" {
		g, err = user.LookupGroup(groupOrId)

		if err != nil {
			g, err = user.LookupGroupId(groupOrId)
		}

		if err != nil {
			gid, err = parseUint(groupOrId)
		}

		if err != nil {
			return
		}
	}

	if u != nil {
		uid, _ = parseUint(u.Uid)
		gid, _ = parseUint(u.Gid)
	}

	if g != nil {
		gid, _ = parseUint(g.Gid)
	}

	return
}

func attrInit(c *exec.Cmd) *syscall.SysProcAttr {
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	return c.SysProcAttr
}

func parseUint(s string) (uint32, error) {
	if out, err := strconv.ParseUint(s, 10, 32); err != nil {
		return 0, err
	} else {
		return uint32(out), nil
	}
}

const (
	SIGKILL = syscall.SIGKILL
	SIGINT  = syscall.SIGINT
)
