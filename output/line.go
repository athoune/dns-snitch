package output

import "fmt"

type Line struct {
	From      string // IP
	Target    string // IP
	TS        int64  `parquet:"timestamp_micros,timestamp(microsecond)"`
	Domain    string
	Port      uint32
	Direction string `parquet:"enum"`
}

func (l *Line) String() string {
	if l.Direction == "UP" {
		return fmt.Sprintf("%v: %s ▶ %s(%s):%d", l.TS, l.From, l.Domain, l.Target, l.Port)
	}
	return fmt.Sprintf("%v: %s ◀ %s(%s):%d", l.TS, l.Target, l.Domain, l.From, l.Port)
}

type LineValue struct {
	Line
	Weight int32
}

type LineValueBySize []*LineValue

func (a LineValueBySize) Len() int           { return len(a) }
func (a LineValueBySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LineValueBySize) Less(i, j int) bool { return a[i].Weight > a[j].Weight }

func Lines2LineValues(k []Line, v []int) []*LineValue {
	values := make([]*LineValue, len(k))
	for i := 0; i < len(k); i++ {
		values[i] = &LineValue{
			Line:   k[i],
			Weight: int32(v[i]),
		}
	}
	return values
}
