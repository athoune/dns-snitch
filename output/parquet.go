package output

import (
	"io"
	"os"
	"time"

	"github.com/parquet-go/parquet-go"
)

var Schema *parquet.Schema

func init() {
	Schema = parquet.SchemaOf(LineValue{})
}

type ParquetWriter struct {
	writer   io.Writer
	duration time.Duration
}

func New(path string) (*ParquetWriter, error) {
	// FIXME handle duration
	w, err := os.OpenFile("./snitch.parquet", os.O_WRONLY+os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}
	return &ParquetWriter{
		writer:   w,
		duration: 10 * time.Second,
	}, nil
}

func NewFromWriter(w io.Writer) *ParquetWriter {
	// FIXME handle duration
	return &ParquetWriter{
		writer:   w,
		duration: 10 * time.Second,
	}
}

func (w *ParquetWriter) Write(k []Line, v []int) error {
	ts := time.Now().Truncate(w.duration).UnixMicro()
	for i := range k {
		k[i].TS = ts
	}

	return parquet.Write[*LineValue](w.writer, Lines2LineValues(k, v))
}
