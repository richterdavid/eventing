package test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/binding/format"
	"github.com/cloudevents/sdk-go/v2/event"
)

// MockStructuredMessage implements a structured-mode message as a simple struct.
// MockStructuredMessage implements both the binding.Message interface and the binding.StructuredWriter
type MockStructuredMessage struct {
	Format format.Format
	Bytes  []byte
}

// Create a new MockStructuredMessage starting from an event.Event. Panics in case of error
func MustCreateMockStructuredMessage(e event.Event) binding.Message {
	return &MockStructuredMessage{
		Bytes:  MustJSON(e),
		Format: format.JSON,
	}
}

func (s *MockStructuredMessage) ReadStructured(ctx context.Context, b binding.StructuredWriter) error {
	return b.SetStructuredEvent(ctx, s.Format, bytes.NewReader(s.Bytes))
}

func (s *MockStructuredMessage) ReadBinary(context.Context, binding.BinaryWriter) error {
	return binding.ErrNotBinary
}

func (bm *MockStructuredMessage) ReadEncoding() binding.Encoding {
	return binding.EncodingStructured
}

func (s *MockStructuredMessage) Finish(error) error { return nil }

func (s *MockStructuredMessage) SetStructuredEvent(ctx context.Context, format format.Format, event io.Reader) (err error) {
	s.Format = format
	s.Bytes, err = ioutil.ReadAll(event)
	if err != nil {
		return
	}

	return nil
}

var _ binding.Message = (*MockStructuredMessage)(nil)
var _ binding.StructuredWriter = (*MockStructuredMessage)(nil)
