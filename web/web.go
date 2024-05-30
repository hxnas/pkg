package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router = chi.Router

func NewMux() Router {
	router := chi.NewMux()
	router.Use(middleware.Recoverer)
	return router
}

// func AddrValidate(addrIn string) (addr, ip, port string, err error) {
// 	addr = addrIn

// 	if len(addrIn) == 0 {
// 		port = strconv.Itoa(int(defaultPort))
// 		addr = ":" + port
// 		return
// 	}

// 	if addrIn[0] == '[' {
// 		n := strings.LastIndexByte(addrIn, ']')
// 		if n == -1 {
// 			err = errors.New("missing ]")
// 			return
// 		}
// 		ip = addrIn[1:n]
// 		if addrIn[n+1] != ':' {
// 			err = errors.New("missing ':'")
// 			return
// 		}
// 		port = addrIn[n+2:]
// 	} else {
// 		ip, port, _ = strings.Cut(addrIn, ":")
// 	}

// 	switch ip {
// 	case "":
// 		ip = "0.0.0.0"
// 	case "localhost", "127.0.0.1", "::1":
// 		ip = "127.0.0.1"
// 	default:
// 		var a netip.Addr
// 		if a, err = netip.ParseAddr(ip); err != nil {
// 			err = errors.New("invalid ip " + strconv.Quote(ip) + " parsing " + strconv.Quote(ip))
// 			return
// 		}
// 		if a.Is6() {
// 			ip = "[" + a.String() + "]"
// 		} else {
// 			ip = a.String()
// 		}
// 	}

// 	switch port {
// 	case "", "0":
// 		port = strconv.FormatUint(uint64(defaultPort), 10)
// 	case "http":
// 		port = "80"
// 	case "https":
// 		port = "443"
// 	default:
// 		var p16 uint64
// 		if p16, err = strconv.ParseUint(port, 10, 16); err != nil {
// 			err = errors.New("invalid port " + strconv.Quote(port) + " parsing " + strconv.Quote(addrIn))
// 			return
// 		}
// 		port = strconv.FormatUint(p16, 10)
// 	}

// 	addr = net.JoinHostPort(ip, port)
// 	return
// }
