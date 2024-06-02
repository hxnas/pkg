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
		recursive := lod.First(recursives)
		if recursive {
			err = os.MkdirAll(path, perm)
		} else {
			err = os.Mkdir(path, perm)
		}
		if err != nil {
			err = lod.Errf("%w", err)
		}
		slog.Log(ctx, lod.ErrDebug(err), "mkdir", "path", path, "perm", perm.String(), "recursive", recursive)
		return
	}
}

func Mkdirs(path string) Caller {
	return Mkdir(path, 0777, true)
}

func ReadDirNames(path string, n int) (files []string, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		err = lod.Errf("%w", err)
		return
	}
	defer f.Close()
	if files, err = f.Readdirnames(n); err != nil {
		err = lod.Errf("%w", err)
	}
	slices.Sort(files)
	return
}

func ReadDirs(path string, n int) (files []fs.DirEntry, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		err = lod.Errf("%w", err)
		return
	}
	defer f.Close()
	if files, err = f.ReadDir(n); err != nil {
		err = lod.Errf("%w", err)
	}
	slices.SortFunc(files, func(a, b fs.DirEntry) int { return cmp.Compare(a.Name(), b.Name()) })
	return
}

func IsSubPath(basepath, targpath string) bool {
	rel, err := filepath.Rel(basepath, targpath)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}
