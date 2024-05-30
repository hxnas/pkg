package sys

import (
	"archive/tar"
	"context"
	"errors"
	"io"
	"io/fs"
)

type TarDecodeFunc func(src io.Reader) (decoded io.ReadCloser, err error)
type TarExtractFunc func(tr *tar.Reader, h *tar.Header) (err error)

func Tar(src io.Reader) *TarFile {
	return &TarFile{src: src}
}

type TarFile struct {
	src        io.Reader
	decodeFunc TarDecodeFunc
}

func (t *TarFile) Decode(decodeFunc TarDecodeFunc) *TarFile {
	t.decodeFunc = decodeFunc
	return t
}

func (t *TarFile) Read(ctx context.Context, extractFunc TarExtractFunc) (err error) {
	var dr io.ReadCloser
	if t.decodeFunc != nil {
		dr, err = t.decodeFunc(t.src)
	} else {
		dr = io.NopCloser(t.src)
	}
	if err != nil {
		return
	}
	defer dr.Close()

	tr := tar.NewReader(dr)
	for h, e := tr.Next(); e != io.EOF; h, e = tr.Next() {
		if e != nil {
			err = e
			return
		}

		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
			err = extractFunc(tr, h)
		}

		if err != nil {
			if errors.Is(err, fs.SkipAll) {
				err = nil
			}
			return
		}
	}
	return
}

func (t *TarFile) Caller(extractFunc TarExtractFunc) Caller {
	return func(ctx context.Context) (err error) { return t.Read(ctx, extractFunc) }
}
