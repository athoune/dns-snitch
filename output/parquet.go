package output

import (
	"io"
	"os"

	"github.com/parquet-go/parquet-go"
)

var Schema *parquet.Schema

func init() {
	Schema = parquet.SchemaOf(LineValue{})
}

type ParquetWriter struct {
	writer io.Writer
}

func New(path string) (*ParquetWriter, error) {
	w, err := os.OpenFile("./snitch.parquet", os.O_WRONLY+os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}
	return &ParquetWriter{w}, nil
}

func NewFromWriter(w io.Writer) *ParquetWriter {
	return &ParquetWriter{w}
}

func (w *ParquetWriter) Write(k []Line, v []int) error {

	return parquet.Write[*LineValue](w.writer, Lines2LineValues(k, v))
}
