package sys

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

var buffer32KPool = &sync.Pool{New: func() interface{} { s := make([]byte, 32*1024); return &s }}

func IOWriteFile(ctx context.Context, src io.Reader, dstPath string, perm fs.FileMode) (err error) {
	if err = os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
		return
	}

	var dst *os.File
	if dst, err = os.Create(dstPath); err == nil {
		if err = IOCopy(ctx, dst, src); err == nil && perm != 0666 {
			err = dst.Chmod(perm)
		}
		if e := dst.Close(); err == nil {
			err = e
		}
	}
	return
}

func IOCopy(ctx context.Context, w io.Writer, r io.Reader) (err error) {
	buf := buffer32KPool.Get().(*[]byte)
	_, err = io.CopyBuffer(IOW(ctx, w), IOR(ctx, r), *buf)
	buffer32KPool.Put(buf)
	return
}

func IORE(r io.Reader, err error) io.ReadCloser {
	if err == nil {
		if rc, ok := r.(io.ReadCloser); ok {
			return rc
		}
		return io.NopCloser(r)
	}
	return io.NopCloser(ioRw(func(p []byte) (int, error) { return 0, err }))
}

func IOW(ctx context.Context, w io.Writer) io.Writer {
	return ioRw(func(p []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return w.Write(p)
		}
	})
}

func IOR(ctx context.Context, r io.Reader) io.Reader {
	return ioRw(func(p []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return r.Read(p)
		}
	})
}

type ioRw func(p []byte) (int, error)

func (fn ioRw) Read(p []byte) (int, error)  { return fn(p) }
func (fn ioRw) Write(p []byte) (int, error) { return fn(p) }
