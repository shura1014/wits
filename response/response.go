package response

import (
	"bufio"
	"net"
	"net/http"
)

type Response interface {
	Unwrap() http.ResponseWriter
	Header() http.Header
	SetHeader(name, value string)
	ContentType() string
	WriteStatus(status int)
	WriteHeader(status int)
	WriteStatusNow()
	Write(b []byte) (n int, err error)
	WriteString(s string) (n int, err error)
	Flush()
	Hijack() (net.Conn, *bufio.ReadWriter, error)
	Status() int
	Size() int64
	Reset(w http.ResponseWriter)
	Buffer() []byte
	BufferString() string
	ClearBuffer()
}
