package sys

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func ReadToFile(ctx context.Context, src io.Reader, dstPath string, perm fs.FileMode) (err error) {
	if err = os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
		return
	}

	var dst *os.File
	if dst, err = os.Create(dstPath); err == nil {
		if err = Copy(ctx, dst, src); err == nil && perm != 0666 {
			err = dst.Chmod(perm)
		}
		if e := dst.Close(); err == nil {
			err = e
		}
	}
	return
}

func Copy(ctx context.Context, w io.Writer, r io.Reader) error {
	return ioRw{Writer: w, Reader: r}.Copy(ctx)
}

func ContextWriter(ctx context.Context, w io.Writer) io.Writer {
	return ioRw{Writer: w, Context: ctx}
}

func ContextReader(ctx context.Context, r io.Reader) io.Writer {
	return ioRw{Reader: r, Context: ctx}
}

func ContextWrap(ctx context.Context, w io.Writer, r io.Reader) io.ReadWriter {
	return ioRw{Reader: r, Writer: w, Context: ctx}
}

var buffer32KPool = &sync.Pool{New: func() interface{} { s := make([]byte, 32*1024); return &s }}

type ioRw struct {
	context.Context
	io.Writer
	io.Reader
}

func (c ioRw) Read(p []byte) (n int, err error) {
	if c.Reader != nil {
		select {
		case <-c.Done():
			return 0, c.Err()
		default:
			return c.Reader.Read(p)
		}
	}
	return 0, errors.New("reader is nil")
}

func (c ioRw) Write(p []byte) (n int, err error) {
	if c.Writer != nil {
		select {
		case <-c.Done():
			return 0, c.Err()
		default:
			return c.Writer.Write(p)
		}
	}
	return 0, errors.New("writer is nil")
}

func (c ioRw) Copy(ctx context.Context) (err error) {
	c.Context = ctx
	buf := buffer32KPool.Get().(*[]byte)
	_, err = io.CopyBuffer(c, c, *buf)
	buffer32KPool.Put(buf)
	return
}
