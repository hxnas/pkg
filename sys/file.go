package sys

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hxnas/pkg/lod"
)

// cat 以utf8编码读取文件文本内容
func Cat(path string) string {
	d, _ := os.ReadFile(path)
	return string(d)
}

// randText 生成一个指定长度的随机字符串。
func RandText(n int) (s string) {
	var d = make([]byte, (n+4)/8*5)
	rand.Read(d)
	if s = B32.EncodeToString(d); len(s) > n {
		s = s[:n]
	}
	return
}

var B32 = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567").WithPadding(base32.NoPadding)

// fileWrite 写入文件，存在则跳过
func FileWrite[T ~string | ~[]byte](path string, data T, perm os.FileMode, overwrite ...bool) Caller {
	return func(ctx context.Context) (err error) {
		if err = Mkdirs(filepath.Dir(path)).Call(ctx); err != nil {
			err = lod.Errf("%w", err)
			return
		}

		err = func() (err error) {
			var f *os.File
			if lod.First(overwrite) {
				f, err = os.Create(path)
			} else {
				f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
			}

			if err != nil {
				if os.IsExist(err) {
					slog.WarnContext(ctx, fmt.Sprintf("%s has exist, ignore write", path))
					return nil
				}
				err = lod.Errf("%w", err)
				return
			}

			if _, err = f.Write([]byte(data)); err != nil {
				err = lod.Errf("%w", err)
			}

			if err == nil && perm > 0 && perm != 0666 {
				if err = f.Chmod(perm); err != nil {
					err = lod.Errf("%w", err)
				}
			}

			if ce := f.Close(); err == nil {
				if err = ce; err != nil {
					err = lod.Errf("%w", err)
				}
			}
			return
		}()

		slog.Log(ctx, lod.ErrDebug(err), "write", "path", path, "data", string(data), "err", err)
		return
	}
}

// fileWriteCopy 先写入文件到 writeTo，存在则跳过写入，再复制到 copyTo，存在则跳过
func FileWriteCopy[T ~string | ~[]byte](writeTo, copyTo string, data T, perm os.FileMode) Caller {
	return func(ctx context.Context) (err error) {
		content := []byte(data)
		var stat os.FileInfo
		if stat, err = os.Stat(writeTo); stat != nil && stat.Mode().IsRegular() {
			content, err = os.ReadFile(writeTo)
		} else if os.IsNotExist(err) {
			err = FileWrite(writeTo, content, 0).Call(ctx)
		}

		if err != nil {
			err = lod.Errf("%w", err)
			return
		}

		return FileWrite(copyTo, content, perm).Call(ctx)
	}
}
