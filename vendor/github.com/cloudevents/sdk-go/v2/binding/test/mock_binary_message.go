package test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/binding/spec"
	"github.com/cloudevents/sdk-go/v2/event"
)

// MockBinaryMessage implements a binary-mode message as a simple struct.
// MockBinaryMessage implements both the binding.Message interface and the binding.BinaryWriter
type MockBinaryMessage struct {
	Metadata   map[spec.Attribute]interface{}
	Extensions map[string]interface{}
	Body       []byte
}

// Create a new MockBinaryMessage starting from an event.Event. Panics in case of error
func MustCreateMockBinaryMessage(e event.Event) binding.Message {
	version := spec.VS.Version(e.SpecVersion())

	m := MockBinaryMessage{
		Metadata:   make(map[spec.Attribute]interface{}),
		Extensions: make(map[string]interface{}),
	}

	for _, attribute := range version.Attributes() {
		val := attribute.Get(e.Context)
		if val != nil {
			m.Metadata[attribute] = val
		}
	}

	for k, v := range e.Extensions() {
		m.Extensions[k] = v
	}

	m.Body = e.Data()

	return &m
}

func (bm *MockBinaryMessage) ReadStructured(context.Context, binding.StructuredWriter) error {
	return binding.ErrNotStructured
}

func (bm *MockBinaryMessage) ReadBinary(ctx context.Context, b binding.BinaryWriter) error {
	err := b.Start(ctx)
	if err != nil {
		return err
	}
	for k, v := range bm.Metadata {
		err = b.SetAttribute(k, v)
		if err != nil {
			return err
		}
	}
	for k, v := range bm.Extensions {
		err = b.SetExtension(k, v)
		if err != nil {
			return err
		}
	}
	if len(bm.Body) != 0 {
		err = b.SetData(bytes.NewReader(bm.Body))
		if err != nil {
			return err
		}
	}
	return b.End(ctx)
}

func (bm *MockBinaryMessage) ReadEncoding() binding.Encoding {
	return binding.EncodingBinary
}

func (bm *MockBinaryMessage) Finish(error) error { return nil }

func (bm *MockBinaryMessage) Start(ctx context.Context) error {
	bm.Metadata = make(map[spec.Attribute]interface{})
	bm.Extensions = make(map[string]interface{})
	return nil
}

func (bm *MockBinaryMessage) SetAttribute(attribute spec.Attribute, value interface{}) error {
	bm.Metadata[attribute] = value
	return nil
}

func (bm *MockBinaryMessage) SetExtension(name string, value interface{}) error {
	bm.Extensions[name] = value
	return nil
}

func (bm *MockBinaryMessage) SetData(data io.Reader) (err error) {
	bm.Body, err = ioutil.ReadAll(data)
	return err
}

func (bm *MockBinaryMessage) End(ctx context.Context) error {
	return nil
}

var _ binding.Message = (*MockBinaryMessage)(nil)
var _ binding.BinaryWriter = (*MockBinaryMessage)(nil)
