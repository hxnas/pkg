package sys

import (
	"archive/tar"
	"context"
	"io"
	"io/fs"

	"github.com/hxnas/pkg/lod"
)

type TarDecoder func(io.Reader) (io.ReadCloser, error)
type TarWalkFunc func(ctx context.Context, r io.Reader, h *tar.Header) (err error)

func TarExtract(ctx context.Context, src io.Reader, walk TarWalkFunc, decoder ...TarDecoder) (err error) {
	dr := lod.Firsts(decoder...).Decode(src)
	defer dr.Close()

	var hdr *tar.Header
	tr := tar.NewReader(dr)
	for hdr, err = tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			return
		}
		err = walk.Read(ctx, tr, hdr)
	}

	if err == fs.SkipAll || err == io.EOF {
		err = nil
	}
	return
}

func (d TarDecoder) Decode(r io.Reader) io.ReadCloser {
	if d != nil {
		if rc, err := d(r); err == nil {
			return rc
		} else {
			return io.NopCloser(ioRw(func(p []byte) (int, error) { return 0, err }))
		}
	}
	return io.NopCloser(r)
}

func (w TarWalkFunc) Read(ctx context.Context, r io.Reader, h *tar.Header) (err error) {
	defer io.Copy(io.Discard, r)
	if w != nil {
		return w(ctx, r, h)
	}
	return nil
}
