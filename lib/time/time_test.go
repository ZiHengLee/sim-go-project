package time

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	n := time.Now()

	var v int64

	v = Unix(n)
	fmt.Printf("unixSecd=%v\n", v)
	if v != n.Unix() {
		t.Errorf("unix unmatch %v:%v", v, n.Unix())
	}

	v = UnixMilli(n)
	fmt.Printf("unixMill=%v\n", v)
	if v != n.UnixNano()/1000000 {
		t.Errorf("unix milli unmatch %v:%v", v, n.Unix())
	}

	v = UnixMicro(n)
	fmt.Printf("unixMicr=%v\n", v)
	if v != n.UnixNano()/1000 {
		t.Errorf("unix micro unmatch %v:%v", v, n.Unix())
	}

	v = UnixNano(n)
	fmt.Printf("unixNano=%v\n", v)
	if v != n.UnixNano() {
		t.Errorf("unix nano unmatch %v:%v", v, n.Unix())
	}
}

func TestTimeFormate(t *testing.T) {
	fmt.Println(DayStrToTimeStamp("2020-04-01", 8, FormatYmd))
}
