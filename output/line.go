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
