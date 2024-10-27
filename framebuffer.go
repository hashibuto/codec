package codec

import (
	"bytes"
	"encoding/binary"
)

type FrameBuffer struct {
	buffer bytes.Buffer

	nextFrameSize uint32
	lengthBuffer  []byte
}

func (fb *FrameBuffer) Write(data []byte) {
	fb.buffer.Write(data)
}

func (fb *FrameBuffer) Read() []byte {
	if fb.buffer.Len() < 4 {
		return nil
	}

	if fb.lengthBuffer == nil {
		fb.lengthBuffer = make([]byte, 4)
	}

	if fb.nextFrameSize == 0 {
		_, _ = fb.buffer.Read(fb.lengthBuffer)
		fb.nextFrameSize = binary.BigEndian.Uint32(fb.lengthBuffer)
	}

	if fb.buffer.Len() < int(fb.nextFrameSize) {
		return nil
	}

	frameBuffer := make([]byte, fb.nextFrameSize)
	_, _ = fb.buffer.Read(frameBuffer)
	fb.nextFrameSize = 0
	return frameBuffer
}

// MakeDataFrame returns a new buffer with the length prepended as a 4 byte big endian unsigned integer.
func MakeDataFrame(data []byte) []byte {
	l := uint32(len(data))
	frameData := make([]byte, len(data)+4)
	binary.BigEndian.PutUint32(frameData, l)
	copy(frameData[4:], data)
	return frameData
}
