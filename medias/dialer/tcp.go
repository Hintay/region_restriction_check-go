package dialer

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net"
)

func TCPDialer(dns string, forceAddr string, logger *log.Entry) fasthttp.DialFunc {
	dialer := fasthttp.TCPDialer{
		Resolver: &net.Resolver{
			PreferGo:     true,
			StrictErrors: false,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				logger.Debugln("connecting to dns")
				return d.DialContext(ctx, network, dns)
			},
		},
	}

	return func(addr string) (net.Conn, error) {
		if len(forceAddr) > 0 {
			addr = forceAddr
		}
		return dialer.Dial(addr)
	}
}
