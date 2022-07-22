package tcp

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"net"
	"time"
)

type Connection struct {
	conn net.Conn

	// Read
	readTimeout time.Duration
	sr          *StreamReader

	// Write
	writeTimeout time.Duration
	sw           *StreamWriter
}

// DialOption specifies an option for dialing a  server.
type DialOption struct {
	f func(*dialOptions)
}

type dialOptions struct {
	readTimeout         time.Duration
	writeTimeout        time.Duration
	tlsHandshakeTimeout time.Duration
	dialer              *net.Dialer
	dialContext         func(ctx context.Context, network, addr string) (net.Conn, error)
	useTLS              bool
	skipVerify          bool
	tlsConfig           *tls.Config
}

// DialTLSHandshakeTimeout specifies the maximum amount of time waiting to
// wait for a TLS handshake. Zero means no timeout.
// If no DialTLSHandshakeTimeout option is specified then the default is 30 seconds.
func DialTLSHandshakeTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.tlsHandshakeTimeout = d
	}}
}

// DialReadTimeout specifies the timeout for reading a single command reply.
func DialReadTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.readTimeout = d
	}}
}

// DialWriteTimeout specifies the timeout for writing a single command.
func DialWriteTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.writeTimeout = d
	}}
}

// DialConnectTimeout specifies the timeout for connecting to the Redis server when
// no DialNetDial option is specified.
// If no DialConnectTimeout option is specified then the default is 30 seconds.
func DialConnectTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.dialer.Timeout = d
	}}
}

// DialKeepAlive specifies the keep-alive period for TCP connections to the Redis server
// when no DialNetDial option is specified.
// If zero, keep-alives are not enabled. If no DialKeepAlive option is specified then
// the default of 5 minutes is used to ensure that half-closed TCP sessions are detected.
func DialKeepAlive(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.dialer.KeepAlive = d
	}}
}

// DialNetDial specifies a custom dial function for creating TCP
// connections, otherwise a net.Dialer customized via the other options is used.
// DialNetDial overrides DialConnectTimeout and DialKeepAlive.
func DialNetDial(dial func(network, addr string) (net.Conn, error)) DialOption {
	return DialOption{func(do *dialOptions) {
		do.dialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dial(network, addr)
		}
	}}
}

// DialContextFunc specifies a custom dial function with context for creating TCP
// connections, otherwise a net.Dialer customized via the other options is used.
// DialContextFunc overrides DialConnectTimeout and DialKeepAlive.
func DialContextFunc(f func(ctx context.Context, network, addr string) (net.Conn, error)) DialOption {
	return DialOption{func(do *dialOptions) {
		do.dialContext = f
	}}
}

// DialTLSConfig specifies the config to use when a TLS connection is dialed.
// Has no effect when not dialing a TLS connection.
func DialTLSConfig(c *tls.Config) DialOption {
	return DialOption{func(do *dialOptions) {
		do.tlsConfig = c
	}}
}

// DialTLSSkipVerify disables server name verification when connecting over
// TLS. Has no effect when not dialing a TLS connection.
func DialTLSSkipVerify(skip bool) DialOption {
	return DialOption{func(do *dialOptions) {
		do.skipVerify = skip
	}}
}

// DialUseTLS specifies whether TLS should be used when connecting to the
// server. This option is ignore by DialURL.
func DialUseTLS(useTLS bool) DialOption {
	return DialOption{func(do *dialOptions) {
		do.useTLS = useTLS
	}}

}

// DialTimeout acts like Dial but takes timeouts for establishing the
// connection to the server, writing a command and reading a reply.
//
// Deprecated: Use Dial with options instead.
func DialTimeout(network, address string, connectTimeout, readTimeout, writeTimeout time.Duration) (*Connection, error) {
	return Dial(network, address,
		DialConnectTimeout(connectTimeout),
		DialReadTimeout(readTimeout),
		DialWriteTimeout(writeTimeout))
}

// Dial connects to the Redis server at the given network and
// address using the specified options.
func Dial(network, address string, options ...DialOption) (*Connection, error) {
	return DialContext(context.Background(), network, address, options...)
}

type tlsHandshakeTimeoutError struct{}

func (tlsHandshakeTimeoutError) Timeout() bool   { return true }
func (tlsHandshakeTimeoutError) Temporary() bool { return true }
func (tlsHandshakeTimeoutError) Error() string   { return "TLS handshake timeout" }

func cloneTLSConfig(cfg *tls.Config) *tls.Config {
	return cfg.Clone()
}

// DialContext connects to the  server at the given network and
// address using the specified options and context.
func DialContext(ctx context.Context, network, address string, options ...DialOption) (*Connection, error) {

	do := dialOptions{
		dialer: &net.Dialer{
			Timeout:   time.Second * 30,
			KeepAlive: time.Minute * 5,
		},
		tlsHandshakeTimeout: time.Second * 10,
	}

	for _, option := range options {
		option.f(&do)
	}

	if do.dialContext == nil {
		do.dialContext = do.dialer.DialContext
	}

	netConn, err := do.dialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}

	if do.useTLS {
		var tlsConfig *tls.Config
		if do.tlsConfig == nil {
			tlsConfig = &tls.Config{InsecureSkipVerify: do.skipVerify}
		} else {
			tlsConfig = cloneTLSConfig(do.tlsConfig)
		}
		if tlsConfig.ServerName == "" {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				netConn.Close()
				return nil, err
			}
			tlsConfig.ServerName = host
		}

		tlsConn := tls.Client(netConn, tlsConfig)
		errc := make(chan error, 2) // buffered so we don't block timeout or Handshake
		if d := do.tlsHandshakeTimeout; d != 0 {
			timer := time.AfterFunc(d, func() {
				errc <- tlsHandshakeTimeoutError{}
			})
			defer timer.Stop()
		}
		go func() {
			errc <- tlsConn.Handshake()
		}()
		if err := <-errc; err != nil {
			// Timeout or Handshake error.
			netConn.Close() // nolint: errcheck
			return nil, err
		}

		netConn = tlsConn
	}

	c := &Connection{
		conn:         netConn,
		sr:           NewStreamReader(netConn),
		sw:           NewStreamWriter(netConn),
		readTimeout:  do.readTimeout,
		writeTimeout: do.writeTimeout,
	}

	return c, nil
}

func (c *Connection) Flush() error {

	if c.writeTimeout != 0 {
		if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return err
		}
	}

	return c.sw.Flush()
}

func (c *Connection) WriteBytes(data []byte) error {

	return c.sw.WriteBytes(data)
}

func (c *Connection) WriteString(data string) error {

	return c.sw.WriteString(data)
}

func (c *Connection) WriteHex(data string) error {

	hexdata, err := hex.DecodeString(data)

	if err != nil {
		return err
	}

	return c.WriteBytes(hexdata)

}

// readLine reads a line of input from the RESP stream.
func (c *Connection) ReadLine() ([]byte, error) {

	var deadline time.Time
	if c.readTimeout != 0 {
		deadline = time.Now().Add(c.readTimeout)
	}

	if err := c.conn.SetReadDeadline(deadline); err != nil {
		return nil, err
	}

	return c.sr.ReadLine()
}

// readLine reads a line of input from the RESP stream.
func (c *Connection) ReadBytes(n int) ([]byte, error) {

	var deadline time.Time
	if c.readTimeout != 0 {
		deadline = time.Now().Add(c.readTimeout)
	}

	if err := c.conn.SetReadDeadline(deadline); err != nil {
		return nil, err
	}

	return c.sr.ReadBytes(n)
}

func (c *Connection) Close() {

	c.conn.Close()
}
