package sys

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hxnas/pkg/lod"
	"github.com/moby/sys/symlink"
)

func RealPath(path string) (string, error) {
	if IsSymlink(path) {
		return symlink.EvalSymlinks(path)
	} else {
		return filepath.Abs(path)
	}
}

func IsSymlink(path string) bool {
	f, err := os.Lstat(path)
	return err == nil && f.Mode()&os.ModeSymlink == os.ModeSymlink
}

// Symlink 建立软链接
func Symlink(srcPath, dstPath string) func(ctx context.Context) (err error) {
	return func(ctx context.Context) (err error) {
		slog.DebugContext(ctx, "links", "src", srcPath, "dst", dstPath)
		if srcPath, err = filepath.Abs(srcPath); err != nil {
			err = lod.Errf("%w", err)
			return
		}

		if dstPath, err = filepath.Abs(dstPath); err != nil {
			err = lod.Errf("%w", err)
			return
		}

		if err = os.MkdirAll(filepath.Dir(dstPath), 0777); err != nil {
			err = lod.Errf("%w", err)
			return
		}

		var dstInfo os.FileInfo
		if dstInfo, err = os.Lstat(dstPath); dstInfo != nil {
			if dstInfo.IsDir() || dstInfo.Mode()&os.ModeSymlink != 0 { //如果是空文件夹或者软链接，删除重建
				err = os.Remove(dstPath)
			} else {
				err = fmt.Errorf("%q is exist, can not make a symlink", dstPath)
			}
		}

		if err != nil && !os.IsNotExist(err) {
			err = lod.Errf("%w", err)
			return
		}

		if err = os.Symlink(srcPath, dstPath); err != nil {
			err = lod.Errf("%w", err)
			return
		}

		return
	}
}
