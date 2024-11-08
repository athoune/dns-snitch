package output

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

type LineValueBySize []*LineValue

func (a LineValueBySize) Len() int           { return len(a) }
func (a LineValueBySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LineValueBySize) Less(i, j int) bool { return a[i].Size > a[j].Size }

func Lines2LineValues(k []*Line, v []int) []*LineValue {
	values := make([]*LineValue, len(k))
	for i := 0; i < len(k); i++ {
		values[i] = &LineValue{
			Line: *k[i],
			Size: int32(v[i]),
		}
	}
	return values
}
