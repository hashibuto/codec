package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrameMaker(t *testing.T) {
	text := "hello world"
	frame := MakeDataFrame([]byte(text))
	assert.Len(t, frame, len(text)+4)
}

func TestFrameBuffer(t *testing.T) {
	text := "hello world"
	frame := MakeDataFrame([]byte(text))

	var fb FrameBuffer

	fb.Write(frame[:5])
	outFrame := fb.Read()
	assert.Nil(t, outFrame)

	fb.Write(frame[5:])
	outFrame = fb.Read()
	assert.Equal(t, []byte(text), outFrame)
}
