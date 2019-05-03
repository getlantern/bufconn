package bufconn

import (
	"bytes"
	"io"
	"testing"

	"github.com/getlantern/mockconn"
	"github.com/stretchr/testify/assert"
)

var (
	alphabet   = "abcdefghijklmnopqrstuvwxyz"
	dataString = alphabet + "\n" + alphabet
	data       = []byte(dataString)
)

func TestRead(t *testing.T) {
	response := bytes.NewBuffer(data)
	wrapped := mockconn.New(&bytes.Buffer{}, response)
	conn := Wrap(wrapped)
	head := conn.Head()
	line1, err := head.ReadString('\n')
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, alphabet+"\n", line1) {
		return
	}
	response.WriteString(alphabet)
	_line2, err := readAll(conn)
	if !assert.NoError(t, err) {
		return
	}
	line2 := string(_line2)
	assert.Equal(t, alphabet+alphabet, line2)
}

func readAll(r io.Reader) ([]byte, error) {
	var result []byte
	b := make([]byte, 1)
	for {
		n, err := r.Read(b)
		if n > 0 {
			result = append(result, b[:n]...)
		}
		if err != nil {
			if err == io.EOF {
				return result, nil
			}
			return nil, err
		}
	}
}
