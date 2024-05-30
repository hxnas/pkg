package sys

import (
	"cmp"
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hxnas/pkg/lod"
)

func Mkdir(path string, perm fs.FileMode, recursives ...bool) Caller {
	return func(ctx context.Context) (err error) {
		recursive := lod.Select(recursives...)
		if recursive {
			err = os.MkdirAll(path, perm)
		} else {
			err = os.Mkdir(path, perm)
		}
		slog.Log(ctx, lod.ErrDebug(err), "mkdir", "path", path, "perm", perm.String(), "recursive", recursive)
		return
	}
}

func Mkdirs(path string) Caller {
	return Mkdir(path, 0777, true)
}

func ReadDirNames(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dirs, err := f.Readdirnames(n)
	slices.Sort(dirs)
	return dirs, err
}

func ReadDirs(path string, n int) ([]fs.DirEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dirs, err := f.ReadDir(n)
	slices.SortFunc(dirs, func(a, b fs.DirEntry) int { return cmp.Compare(a.Name(), b.Name()) })
	return dirs, err
}

func IsSubPath(basepath, targpath string) bool {
	rel, err := filepath.Rel(basepath, targpath)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}
