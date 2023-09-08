package GoWinHttp

import (
	"net"
)

type MyConn struct {
	net.Conn
	hook func(b []byte)
}

func (c *MyConn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	c.hook(b[0:n])

	/*
		if c.s == nil {
			return n, err
		}
		if len(*(c.s)) > 0 {
			return n, err
		}
		c.b.Write(b[0:n])
		m := c.b.Bytes()
		y := bytes.Index(m, []byte("\r\n\r\n"))
		if y != -1 {
			h := make(http.Header)
			arr := strings.Split(string(CopyBytes(m[0:y])), "\r\n")
			c.b.Reset()
			for _, v := range arr {
				arr2 := strings.Split(v, ":")
				if len(arr2) >= 1 {
					if len(v) >= len(arr2[0])+1 {
						da := strings.TrimSpace(v[len(arr2[0])+1:])
						if len(h[arr2[0]]) > 0 {
							h[arr2[0]] = append(h[arr2[0]], da)
						} else {
							h[arr2[0]] = []string{da}
						}
					}
				}
			}
			*c.s = h
		}
	*/
	return n, err
}

func (c *MyConn) Write(b []byte) (int, error) {
	return c.Conn.Write(b)
}

func (c *MyConn) Close() error {
	if c.Conn == nil {
		return nil
	}
	return c.Conn.Close()
}
