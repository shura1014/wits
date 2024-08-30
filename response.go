package wits

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
)

type Response struct {
	Writer http.ResponseWriter
	status int
	// 是否已经写入头部
	isWrite bool
	// RestResult 包装，只包装一次
	isRestWrap bool
	size       int64
	Committed  bool
	// 是否被劫持
	hijacked bool
	buffer   *bytes.Buffer

	// 兼容proxy
	// 如果已经flush过了就不再flush
	isFlush bool
}

func Wrap(w http.ResponseWriter) (r *Response) {
	return &Response{
		Writer: w,
		status: http.StatusOK,
		buffer: bytes.NewBuffer(nil),
	}
}

func (r *Response) Unwrap() http.ResponseWriter {
	return r.Writer
}

func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

func (r *Response) SetHeader(name, value string) {
	r.Writer.Header().Set(name, value)
}

func (r *Response) ContentType() string {
	return r.Writer.Header().Get("Content-Type")
}

func (r *Response) WriteStatus(status int) {
	if r.Committed {
		Warn("Headers were already written")
		return
	}
	r.status = status
	r.Committed = true
}

func (r *Response) WriteHeader(status int) {
	r.WriteStatus(status)
}

func (r *Response) WriteStatusNow() {
	if !r.isWrite {
		r.Writer.WriteHeader(r.status)
		r.isWrite = true
	}
}

func (r *Response) Write(b []byte) (n int, err error) {
	n, err = r.buffer.Write(b)
	//n, err = r.Writer.Write(b)
	r.size += int64(n)
	return
}

func (r *Response) WriteString(s string) (n int, err error) {
	r.buffer.WriteString(s)
	//n, err = io.WriteString(r.Writer, s)
	r.size += int64(n)
	return
}

func (r *Response) Flush() {

	if r.hijacked {
		return
	}

	status := r.status
	if status == restErrorCode {
		r.status = http.StatusOK
	}

	r.WriteStatusNow()
	if r.buffer.Len() > 0 {
		_, _ = r.Writer.Write(r.buffer.Bytes())
	}
	r.buffer.Reset()
	r.Writer.(http.Flusher).Flush()
	r.isFlush = true

}

func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	r.hijacked = true
	return r.Writer.(http.Hijacker).Hijack()
}

func (r *Response) Status() int {
	return r.status
}
func (r *Response) Size() int64 {
	return r.size
}

func (r *Response) Reset(w http.ResponseWriter) {
	r.Writer = w
	r.isWrite = false
	r.isRestWrap = false
	r.Committed = false
	r.isFlush = false
	r.size = 0
	r.status = http.StatusOK
	r.buffer.Reset()
}

func (r *Response) Buffer() []byte {
	return r.buffer.Bytes()
}

func (r *Response) BufferString() string {
	return r.buffer.String()
}

func (r *Response) ClearBuffer() {
	r.buffer.Reset()
}
