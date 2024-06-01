//go:build !windows && !darwin
// +build !windows,!darwin

package sys

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/hxnas/pkg/lod"
	"golang.org/x/sys/unix"
)

func Chroot(target string) Caller {
	return func(ctx context.Context) (err error) {
		defer func() { slog.Log(ctx, lod.ErrDebug(err), "chroot", "target", target, "err", err) }()
		if err = os.MkdirAll(target, 0777); err != nil {
			return
		}
		if err = unix.Chdir(target); err != nil {
			return
		}
		if err = unix.Chroot("."); err == nil {
			return
		}
		return
	}
}

func Chdir(target string) Caller {
	return func(ctx context.Context) (err error) {
		err = unix.Chdir(target)
		slog.Log(ctx, lod.ErrDebug(err), "chdir", "target", target, "err", err)
		return
	}
}

func Chown(path string, uid, gid uint32, recursive ...bool) Caller {
	return func(ctx context.Context) (err error) {
		return fileWalk(path, func(cur string) (err error) {
			err = unix.Chown(cur, int(uid), int(gid))
			slog.Log(ctx, lod.ErrDebug(err), "chown", "uid", uid, "gid", gid, "path", cur, "err", err)
			return
		}, recursive...)
	}
}

func Chmod(path string, perm fs.FileMode, recursive ...bool) Caller {
	return func(ctx context.Context) (err error) {
		return fileWalk(path, func(cur string) (err error) {
			err = unix.Chmod(cur, uint32(perm))
			slog.Log(ctx, lod.ErrDebug(err), "chmod", "perm", perm.String(), "path", cur, "err", err)
			return
		}, recursive...)
	}
}

func Chtimes(name string, atime, mtime, ctime time.Time) Caller {
	return func(ctx context.Context) (err error) {
		err = os.Chtimes(name, atime, mtime)
		slog.Log(ctx, lod.ErrDebug(err), "chtimes", "atime", atime.String(), "mtime", mtime, "ctime", ctime, "err", err)
		return
	}
}

func fileWalk(path string, do func(path string) error, recursive ...bool) error {
	if lod.First(recursive) {
		return filepath.WalkDir(path, func(p string, _ fs.DirEntry, _ error) (err error) { return do(p) })
	} else {
		return do(path)
	}
}
