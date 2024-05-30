package log

import (
	"bufio"
	"context"
	"io"
	"log"
)

func NewWriter(ctx context.Context, handleMessage func(ctx context.Context, msg string)) *Writer {
	lw := &Writer{}
	lw.ctx, lw.cancel = context.WithCancel(ctx)
	lw.pr, lw.pw = io.Pipe()
	lw.done = make(chan struct{})

	go func() {
		defer close(lw.done)
		defer lw.pr.Close()
		defer lw.pw.Close()

		select {
		case <-ctx.Done():
			return
		default:
			lw.Scan()
		}
	}()

	return lw
}

type Writer struct {
	pr            io.ReadCloser
	pw            io.WriteCloser
	ctx           context.Context
	cancel        context.CancelFunc
	handleMessage func(ctx context.Context, msg string)
	done          chan struct{}
}

func (lw *Writer) Scan() {
	for br := bufio.NewScanner(lw.pr); br.Scan(); {
		select {
		case <-lw.ctx.Done():
			return
		default:
			lw.handleMessage(lw.ctx, br.Text())
		}
	}
}

func (lw *Writer) Write(p []byte) (n int, err error) {
	n, err = lw.pw.Write(p)
	if err != nil {
		lw.cancel()
	}
	return
}

func (lw *Writer) Close() (err error) {
	if lw.cancel != nil {
		lw.cancel()
	}
	<-lw.done
	return
}

func (lw *Writer) StandardLogger(prefix string) *log.Logger {
	return log.New(lw, prefix, 0)
}
