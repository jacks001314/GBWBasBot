package tcp

import (
	"bufio"
	"net"
)

type StreamWriter struct {
	bw *bufio.Writer
}

func NewStreamWriter(conn net.Conn) *StreamWriter {

	return &StreamWriter{bw: bufio.NewWriter(conn)}

}

func (sw *StreamWriter) Flush() error {

	if err := sw.bw.Flush(); err != nil {
		return err
	}

	return nil
}

func (sw *StreamWriter) WriteBytes(data []byte) (err error) {

	_, err = sw.bw.Write(data)
	return
}

func (sw *StreamWriter) WriteString(data string) (err error) {

	_, err = sw.bw.WriteString(data)

	return
}
