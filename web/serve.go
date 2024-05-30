package web

import (
	"context"
	"log/slog"
	"net"
	"net/http"
)

type ServeOption struct {
	Handler    http.Handler
	Addr       string
	OnShutDown func(ctx context.Context)
}

func Serve(ctx context.Context, options ServeOption) (port int, err error) {
	portc := make(chan int)
	errc := make(chan error)

	s := &http.Server{
		Addr:    options.Addr,
		Handler: options.Handler,
		BaseContext: func(l net.Listener) context.Context {
			addr, _ := l.Addr().(*net.TCPAddr)
			if addr != nil {
				trySend(portc, addr.Port)
			}
			slog.InfoContext(ctx, "web started", "listen", l.Addr().String())
			return ctx
		},
	}

	onShutDown := func() {
		if options.OnShutDown != nil {
			options.OnShutDown(ctx)
		}
	}

	done := make(chan struct{})

	go func() {
		select {
		case <-done:
			return
		case <-ctx.Done():
			_ = s.Shutdown(context.Background())
		}
	}()

	go func() {
		defer close(done)
		defer close(portc)
		defer onShutDown()

		err := s.ListenAndServe()
		trySend(errc, err)
		if err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "web is done!", "err", err)
			return
		}
		slog.InfoContext(ctx, "web is done!")
	}()

	select {
	case port = <-portc:
		close(errc)
		err = <-errc
	case err = <-errc:
		close(errc)
	}
	return
}

func trySend[T any](c chan T, v T) {
	select {
	case <-c:
	default:
		c <- v
	}
}
