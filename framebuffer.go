package codec

import (
	"bytes"
	"encoding/binary"
)

type Config struct {
	LengthOffset         int
	LengthFieldInBytes   int
	LengthFieldByteOrder binary.ByteOrder
	SizeBiasInBytes      int // The size bias is the number of bytes which should be added (including negative numbers) to the length in order to determine the entire frame size

	minBufferSize int
}

func Decode1Byte(byteOrder binary.ByteOrder, source []byte) int {
	return int(source[0])
}

func Decode2Byte(byteOrder binary.ByteOrder, source []byte) int {
	return int(byteOrder.Uint16(source))
}

func Decode4Byte(byteOrder binary.ByteOrder, source []byte) int {
	return int(byteOrder.Uint32(source))
}

func Decode8Byte(byteOrder binary.ByteOrder, source []byte) int {
	return int(byteOrder.Uint64(source))
}

type FrameBuffer struct {
	config Config

	buffer            bytes.Buffer
	nextFrameSize     int
	lengthBuffer      []byte
	lengthDecoderFunc func(binary.ByteOrder, []byte) int
}

func NewFrameBuffer(config Config) *FrameBuffer {
	if config.LengthFieldByteOrder == nil {
		config.LengthFieldByteOrder = binary.BigEndian
	}

	if config.LengthFieldInBytes == 0 {
		config.LengthFieldInBytes = 4
	}

	config.minBufferSize = config.LengthFieldInBytes + config.LengthOffset

	var decodeLength func(binary.ByteOrder, []byte) int

	switch config.LengthFieldInBytes {
	case 1:
		decodeLength = Decode1Byte
	case 2:
		decodeLength = Decode2Byte
	case 4:
		decodeLength = Decode4Byte
	case 8:
		decodeLength = Decode8Byte
	}

	return &FrameBuffer{
		config:            config,
		lengthBuffer:      make([]byte, config.minBufferSize),
		lengthDecoderFunc: decodeLength,
	}
}

func (fb *FrameBuffer) Write(data []byte) {
	fb.buffer.Write(data)
}

func (fb *FrameBuffer) Read() []byte {
	if fb.buffer.Len() < fb.config.minBufferSize {
		return nil
	}

	if fb.nextFrameSize == 0 {
		_, _ = fb.buffer.Read(fb.lengthBuffer)
		fb.nextFrameSize = fb.lengthDecoderFunc(fb.config.LengthFieldByteOrder, fb.lengthBuffer)
	}

	if fb.buffer.Len() < fb.nextFrameSize-fb.config.minBufferSize {
		return nil
	}

	frameBuffer := make([]byte, fb.nextFrameSize+fb.config.SizeBiasInBytes)
	copy(frameBuffer, fb.lengthBuffer)
	_, _ = fb.buffer.Read(frameBuffer[fb.config.minBufferSize:])
	fb.nextFrameSize = 0
	return frameBuffer
}

// MakeSimpleDataFrame returns a new buffer with the length prepended as a 4 byte big endian unsigned integer.
// The length includes its own size in the calculation
func MakeSimpleDataFrame(data []byte) []byte {
	l := uint32(len(data) + 4)
	frameData := make([]byte, len(data)+4)
	binary.BigEndian.PutUint32(frameData, l)
	copy(frameData[4:], data)
	return frameData
}
