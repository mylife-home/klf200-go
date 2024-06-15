package transport

import (
	"bytes"
	"fmt"
	"io"
)

const frameEnd byte = 192
const frameEsc byte = 219
const frameEscEnd byte = 220
const frameEscEsc byte = 221

func SlipEncode(data *bytes.Buffer) *bytes.Buffer {
	out := &bytes.Buffer{}
	out.Grow(data.Len() + 2)

	out.WriteByte(frameEnd)

	for {

		b, err := data.ReadByte()

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		switch b {
		case frameEnd:
			out.WriteByte(frameEsc)
			out.WriteByte(frameEscEnd)

		case frameEsc:
			out.WriteByte(frameEsc)
			out.WriteByte(frameEscEsc)

		default:
			out.WriteByte(b)
		}
	}

	out.WriteByte(frameEnd)

	return out
}

type SlipDecoder struct {
	frames       []*bytes.Buffer
	currentFrame *bytes.Buffer
	escaping     bool
}

func (decoder *SlipDecoder) Reset() {
	decoder.frames = nil
	decoder.currentFrame = nil
	decoder.escaping = false
}

func (decoder *SlipDecoder) NextFrame() *bytes.Buffer {
	if decoder.frames == nil {
		return nil
	}

	var out *bytes.Buffer

	out, decoder.frames = decoder.frames[0], decoder.frames[1:]

	if len(decoder.frames) == 0 {
		decoder.frames = nil
	}

	return out
}

func (decoder *SlipDecoder) AddRaw(data []byte) error {
	for _, b := range data {
		if err := decoder.decode(b); err != nil {
			return err
		}
	}

	return nil
}

func (decoder *SlipDecoder) decode(b byte) error {
	if decoder.currentFrame == nil {
		if b != frameEnd {
			return fmt.Errorf("expected frame start")
		}

		decoder.currentFrame = &bytes.Buffer{}
		decoder.currentFrame.Grow(256) // payload is 25 max, + header/footer
		return nil
	}

	if decoder.escaping {
		switch b {
		case frameEscEnd:
			decoder.currentFrame.WriteByte(frameEnd)
		case frameEscEsc:
			decoder.currentFrame.WriteByte(frameEsc)
		default:
			return fmt.Errorf("bad escape sequence")
		}

		decoder.escaping = false
		return nil
	}

	switch b {
	case frameEnd:
		decoder.frames = append(decoder.frames, decoder.currentFrame)
		decoder.currentFrame = nil

	case frameEsc:
		decoder.escaping = true

	default:
		decoder.currentFrame.WriteByte(b)
	}

	return nil
}
