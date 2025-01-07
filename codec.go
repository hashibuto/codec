package codec

// Codec is an amalgamation of the terms code/decode, indicating that it has the capacity to both encode and decode.  This module
// provides an interface for block level encoding/decoding.  The implication of this block level interface, is that neither encode,
// nor decode will be passed partial data.

type Codec interface {
	// Encodes a byte array and returns the encoded byte array.  If no encoded data is available, `nil` will be returned.
	// If an empty/null payload is passed to encode, the chain of codecs will be flushed for any remaining data which has not
	// yet been encoded.
	Encode(data []byte) ([]byte, error)

	// Decodes a byte array and returns the decoded byte array.  If a decoder buffers and/or insufficient data is
	// available, `nil` will be returned.
	Decode([]byte) ([]byte, error)
}

// CodecChain returns a Codec interface representing a chain of codecs
type CodecChain struct {
	codecs []Codec
}

func NewCodecChain(codecs ...Codec) *CodecChain {
	cChain := &CodecChain{}
	cChain.codecs = append(cChain.codecs, codecs...)

	return cChain
}

func (cc *CodecChain) Encode(data []byte) ([]byte, error) {
	isFlush := len(data) == 0
	var curData = data
	var err error
	for _, codec := range cc.codecs {
		curData, err = codec.Encode(curData)
		if err != nil {
			return nil, err
		}
		if curData == nil && !isFlush {
			return nil, nil
		}
	}

	return curData, nil
}

func (cc *CodecChain) Decode(data []byte) ([]byte, error) {
	isFlush := len(data) == 0
	var curData = data
	var err error
	for i := 0; i < len(cc.codecs); i++ {
		codec := cc.codecs[len(cc.codecs)-1-i]
		curData, err = codec.Decode(curData)
		if err != nil {
			return nil, err
		}
		if curData == nil && !isFlush {
			return nil, nil
		}
	}

	return curData, nil
}
