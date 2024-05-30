//go:build !darwin && !windows
// +build !darwin,!windows

package sys

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hxnas/pkg/lod"
	"github.com/moby/sys/mount"
	"github.com/moby/sys/mountinfo"
)

func Mount(device, target, mType, options string, force ...bool) Caller {
	return func(ctx context.Context) (err error) {
		if !lod.Select(force...) {
			var mounted bool
			if mounted, err = mountinfo.Mounted(target); err != nil || mounted {
				return
			}
		}
		if err = os.MkdirAll(target, 0777); err != nil {
			return
		}

		err = mount.Mount(device, target, mType, options)
		slog.DebugContext(ctx, "mount", "device", device, "target", target, "type", mType, "options", options, "err", err)

		return
	}
}

func Unmount(target string, recursives ...bool) Caller {
	return func(ctx context.Context) (err error) {
		recursive := lod.Select(recursives...)
		if recursive {
			err = mount.RecursiveUnmount(target)
		} else {
			err = mount.Unmount(target)
		}

		slog.DebugContext(ctx, "unmount", "target", target, "recursive", recursive, "err", err)
		return
	}
}

func Bind(srcPath, rootDir string, recursive ...bool) Caller {
	return Mount(srcPath, filepath.Join(rootDir, srcPath), "none", lod.Iif(lod.Select(recursive...), "rbind", "bind"))
}
