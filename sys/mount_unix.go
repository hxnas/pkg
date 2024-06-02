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
		if !lod.First(force) {
			var mounted bool
			if mounted, err = mountinfo.Mounted(target); err != nil || mounted {
				if err != nil {
					if os.IsNotExist(err) {
						err = nil
					} else {
						err = lod.Errf("%w", err)
					}
				}
				return
			}
		}
		if err = os.MkdirAll(target, 0777); err != nil {
			err = lod.Errf("%w", err)
			return
		}

		err = mount.Mount(device, target, mType, options)
		if err != nil {
			err = lod.Errf("%w", err)
		}

		slog.DebugContext(ctx, "mount", "device", device, "target", target, "type", mType, "options", options, "err", err)
		return
	}
}

func Unmount(target string, recursives ...bool) Caller {
	return func(ctx context.Context) (err error) {
		recursive := lod.First(recursives)
		if recursive {
			err = mount.RecursiveUnmount(target)
		} else {
			err = mount.Unmount(target)
		}
		if err != nil {
			err = lod.Errf("%w", err)
		}

		slog.DebugContext(ctx, "unmount", "target", target, "recursive", recursive, "err", err)
		return
	}
}

func Bind(srcPath, rootDir string, recursive ...bool) Caller {
	return Mount(srcPath, filepath.Join(rootDir, srcPath), "none", lod.Iif(lod.First(recursive), "rbind", "bind"))
}
