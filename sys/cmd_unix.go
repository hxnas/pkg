//go:build !darwin && !windows
// +build !darwin,!windows

package sys

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"

	"github.com/hxnas/pkg/lod"
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

const keyTag = "GO_CMD_FORK_TAG"

func IsForkTag(tag string) bool {
	return os.Getenv(keyTag) == tag
}

func Fork(args []string, env []string, tag string) Caller {
	return func(ctx context.Context) (err error) {
		var executable string

		if executable, err = os.Executable(); err != nil {
			slog.WarnContext(ctx, tag+" run", "err", err)
			return
		}

		cmd := exec.CommandContext(ctx, executable, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = NewEnv().Append(env...).Set(keyTag, tag).Environ()

		SetNewNS(cmd)
		SetPdeathsig(cmd, SIGINT)

		slog.DebugContext(ctx, tag+" run", "command", cmd.String())
		if err = cmd.Start(); err != nil {
			err = lod.Errf("%w", err)
			slog.WarnContext(ctx, tag+" run", "err", err)
		} else {
			slog.DebugContext(ctx, tag+" run", "pid", cmd.Process.Pid)
		}

		return cmd.Wait()
	}
}
