package tcp

import (
	"fmt"
	"net"
	"time"
)

//
type Conn struct {
	conn interface {
		Read(b []byte) (n int, err error)
		Write(b []byte) (n int, err error)
		Close() error
		LocalAddr() net.Addr
		RemoteAddr() net.Addr
	}
	dataReader *internalDataReader
	dataWriter *internalDataWriter
}

// Connection to TCP address:port
func Connection(ip, port, user, pass string) (*Conn, error) {
	//fmt.Printf("Trying %s:%d...\n", host, p)
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
	if nil != err {
		return nil, err
	}
	dataReader := newDataReader(conn)
	dataWriter := newDataWriter(conn)
	clientConn := Conn{
		conn:       conn,
		dataReader: dataReader,
		dataWriter: dataWriter,
	}
	return &clientConn, nil
}

// Close TCP Connection
func (clientConn *Conn) Close() error {
	return clientConn.conn.Close()
}

// Read from TCP Connection
func (clientConn *Conn) Read(p []byte) (n int, err error) {
	return clientConn.dataReader.Read(p)
}

// Write to TCP Connection
func (clientConn *Conn) Write(p []byte) (n int, err error) {
	return clientConn.dataWriter.Write(p)
}

// LocalAddr Connection
func (clientConn *Conn) LocalAddr() net.Addr {
	return clientConn.conn.LocalAddr()
}

// RemoteAddr Connection
func (clientConn *Conn) RemoteAddr() net.Addr {
	return clientConn.conn.RemoteAddr()
}
