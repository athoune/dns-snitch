package output

import (
	"io"
	"os"

	"github.com/parquet-go/parquet-go"
)

var Schema *parquet.Schema

type Line struct {
	From      string // IP
	Target    string // IP
	TS        int64  `parquet:"timestamp_micros,timestamp(microsecond)"`
	Domain    string
	Port      uint32
	Direction string `parquet:"enum"`
}

type LineValue struct {
	Line
	Size int32
}

func init() {
	Schema = parquet.SchemaOf(LineValue{})
}

type Writer struct {
	w io.Writer
}

func New(path string) (*Writer, error) {
	w, err := os.OpenFile("./snitch.parquet", os.O_WRONLY+os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}
	return &Writer{w}, nil
}

func NewFromWriter(w io.Writer) *Writer {
	return &Writer{w}
}

func (w *Writer) Write(k []*Line, v []int) error {
	values := make([]*LineValue, len(k))
	for i := 0; i < len(k); i++ {
		values[i] = &LineValue{
			Line: *k[i],
			Size: int32(v[i]),
		}
	}
	return parquet.Write[*LineValue](w.w, values)
}
