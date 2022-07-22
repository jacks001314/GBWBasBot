package tcp

import (
	"bufio"
	"net"
)

type StreamReader struct {
	br *bufio.Reader
}

func NewStreamReader(conn net.Conn) *StreamReader {

	return &StreamReader{br: bufio.NewReader(conn)}
}

// readLine reads a line of input from the RESP stream.
func (sr *StreamReader) ReadLine() ([]byte, error) {

	// To avoid allocations, attempt to read the line using ReadSlice. This
	// call typically succeeds. The known case where the call fails is when
	// reading the output from the MONITOR command.
	p, err := sr.br.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		// The line does not fit in the bufio.Reader's buffer. Fall back to
		// allocating a buffer for the line.
		buf := append([]byte{}, p...)
		for err == bufio.ErrBufferFull {
			p, err = sr.br.ReadSlice('\n')
			buf = append(buf, p...)
		}
		p = buf
	}
	if err != nil {
		return nil, err
	}

	i := len(p) - 2
	if i < 0 || p[i] != '\r' {
		return p, nil
	}
	return p[:i], nil
}

func (sr *StreamReader) ReadBytes(n int) ([]byte, error) {

	data := make([]byte, n)

	_, err := sr.br.Read(data)

	return data, err
}
