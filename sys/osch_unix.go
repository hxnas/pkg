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
		if err = unix.Chdir(target); err == nil {
			err = unix.Chroot(".")
		}
		slog.Log(ctx, lod.ErrDebug(err), "chroot", "target", target, "err", err)
		return
	}
}

func ChrootRun(target string, run Caller) Caller {
	return func(ctx context.Context) (err error) {
		slog.DebugContext(ctx, "chroot run", "target", target)
		var wd string
		if wd, err = os.Getwd(); err != nil {
			return
		}
		slog.DebugContext(ctx, "chroot run", "wd", wd)

		var f *os.File
		if f, err = os.Open(wd); err != nil {
			return
		}
		defer f.Close()

		err = run.Call(ctx)

		slog.DebugContext(ctx, "chroot run", "chdir", wd)
		if e := f.Chdir(); e != nil && err == nil {
			err = e
			return
		}

		slog.DebugContext(ctx, "chroot run", "chroot", ".")
		if e := unix.Chroot("."); e != nil && err == nil {
			err = e
		}

		return
	}
}

func Chdir(target string) Caller {
	return func(ctx context.Context) (err error) {
		err = unix.Chdir(target)
		slog.Log(ctx, lod.ErrDebug(err), "chroot", "target", target, "err", err)
		return
	}
}

func Chown(path string, uid, gid uint32, recursive ...bool) Caller {
	return func(ctx context.Context) (err error) {
		fileWalk(path, lod.Select(recursive...), func(cur string) {
			err := unix.Chown(cur, int(uid), int(gid))
			slog.Log(ctx, lod.ErrDebug(err), "chown", "uid", uid, "gid", gid, "path", cur, "err", err)
		})
		return
	}
}

func Chmod(path string, perm fs.FileMode, recursive ...bool) Caller {
	return func(ctx context.Context) (err error) {
		fileWalk(path, lod.Select(recursive...), func(cur string) {
			err := unix.Chmod(cur, uint32(perm))
			slog.Log(ctx, lod.ErrDebug(err), "chmod", "perm", perm.String(), "path", cur, "err", err)
		})
		return
	}
}

func Chtime(name string, atime, mtime, ctime time.Time) Caller {
	return func(ctx context.Context) (err error) {
		err = os.Chtimes(name, atime, mtime)
		slog.Log(ctx, lod.ErrDebug(err), "chmod", "atime", atime.String(), "mtime", mtime, "ctime", ctime, "err", err)
		return
	}
}

func fileWalk(path string, recursive bool, do func(path string)) {
	if recursive {
		filepath.Walk(path, func(p string, _ fs.FileInfo, _ error) (err error) { do(p); return })
	} else {
		do(path)
	}
}
