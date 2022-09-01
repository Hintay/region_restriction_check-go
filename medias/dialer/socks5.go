package dialer

import (
	"net"
	"net/url"

	"github.com/valyala/fasthttp"
	"golang.org/x/net/proxy"
)

func SocksDialer(socks string, forceAddr string) fasthttp.DialFunc {
	var (
		u      *url.URL
		err    error
		dialer proxy.Dialer
	)
	if u, err = url.Parse(socks); err == nil {
		dialer, err = proxy.FromURL(u, proxy.Direct)
	}
	// It would be nice if we could return the error here. But we can't
	// change our API so just keep returning it in the returned Dial function.
	// Besides, the implementation of proxy.SOCKS5() at the time of writing this
	// will always return nil as error.

	return func(addr string) (net.Conn, error) {
		if err != nil {
			return nil, err
		}
		if len(forceAddr) > 0 {
			addr = forceAddr
		}
		return dialer.Dial("tcp", addr)
	}
}
