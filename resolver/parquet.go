package resolver

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/parquet-go/parquet-go"
)

var Schema *parquet.Schema

type Line struct {
	From      []byte // IP
	Target    []byte // IP
	TS        int64  `parquet:"timestamp_micros,timestamp(microsecond)"`
	Size      uint32
	Domain    string
	Port      uint32
	Direction string `parquet:"enum"`
}

func init() {
	Schema = parquet.SchemaOf(Line{})
}

func (r *Resolver) parquetLoop(errorChan chan error) error {
	w, err := os.OpenFile("./snitch.parquet", os.O_WRONLY+os.O_CREATE, 0660)
	if err != nil {
		errorChan <- err
		return err
	}
	count := make(chan interface{})
	ticker := time.NewTicker(10 * time.Second)
	mutex := &sync.Mutex{}
	lines := make([]*Line, 0)

	writeParquet := func() {
		mutex.Lock()
		if len(lines) > 0 {
			parquet.Write[*Line](w, lines)
			fmt.Println("Wrote Parquet", len(lines))
			lines = make([]*Line, 0)
		}
		mutex.Unlock()
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				writeParquet()
			case <-count:
				writeParquet()
			}
		}
	}()
	for {
		mutex.Lock()
		lines = append(lines, <-r.lines)
		mutex.Unlock()
		if len(lines) >= 1000 {
			count <- new(interface{})
		}
	}
}
