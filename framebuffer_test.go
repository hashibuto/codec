package codec

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrameMaker(t *testing.T) {
	text := "hello world"
	frame := MakeSimpleDataFrame([]byte(text))
	assert.Len(t, frame, len(text)+4)
}

func TestFrameBuffer(t *testing.T) {
	text := "hello world"
	frame := MakeSimpleDataFrame([]byte(text))

	fb := NewFrameBuffer(Config{
		LengthFieldInBytes:   4,
		LengthFieldByteOrder: binary.BigEndian,
	})

	fb.Write(frame[:5])
	outFrame := fb.Read()
	assert.Nil(t, outFrame)

	fb.Write(frame[5:])
	outFrame = fb.Read()
	assert.Equal(t, frame, outFrame)
}
