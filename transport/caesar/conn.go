package caesar

import (
	"net"
	"wwww/transport"
)

type CaesarConn struct {
	conn  transport.TransportConn
	shift int

	debugHook func(bytein, byteout []byte, msg string)
}

func (c *CaesarConn) Read(b []byte) (n int, err error) {
	n, err = c.conn.Read(b)
	if err != nil {
		return n, err
	}
	shifted := Decrypt(b[:n], c.shift)
	if c.debugHook != nil {
		c.debugHook(b[:n], shifted, "CaesarConn Read")
	}
	copy(b[:n], shifted)
	return n, nil
}

func (c *CaesarConn) Write(b []byte) (n int, err error) {
	shifted := Encrypt(b, c.shift)
	n, err = c.conn.Write(shifted)
	if err != nil {
		return n, err
	}
	if c.debugHook != nil {
		c.debugHook(b[:n], shifted, "CaesarConn Write")
	}
	return n, err
}

func (c *CaesarConn) Close() error {
	return c.conn.Close()
}

func (c *CaesarConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *CaesarConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func Encrypt(plaintext []byte, shift int) []byte {
	// 简单的凯撒编码，位移量为shift
	ciphertext := make([]byte, len(plaintext))

	for i, b := range plaintext {
		ciphertext[i] = byte((int(b) + shift) % 256)
	}

	return ciphertext
}

func Decrypt(ciphertext []byte, shift int) []byte {
	// 简单的凯撒解码，位移量为shift
	plaintext := make([]byte, len(ciphertext))

	for i, b := range ciphertext {
		plaintext[i] = byte((int(b) - shift + 256) % 256)
	}

	return plaintext
}
