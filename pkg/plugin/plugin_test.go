package plugin

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func TestAlignByPeriod(t *testing.T) {
	frame := data.NewFrame("foo",
		data.NewField("time", nil, []time.Time{
			time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
		}),
		data.NewField("series1", nil, []float64{
			0,
			10,
			20,
		}),
		data.NewField("series2", nil, []float64{
			20,
			10,
			0,
		}),
	)

	got, err := alignByPeriod(frame, "day")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(frame, got) {
		t.Fatal("should be equal")
	}
}

func TestAlignByPeriod_Month(t *testing.T) {
	frame := data.NewFrame("foo",
		data.NewField("time", nil, []time.Time{
			time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
		}),
		data.NewField("series1", nil, []float64{
			0,
			10,
			20,
		}),
		data.NewField("series2", nil, []float64{
			20,
			10,
			0,
		}),
	)

	got, err := alignByPeriod(frame, "month")
	if err != nil {
		t.Fatal(err)
	}

	if got.Rows() != 1 {
		t.Fatalf("unexpected number of rows: want = %d; got = %d", 1, got.Rows())
	}

	for i := 0; i < got.Rows(); i++ {
		for _, f := range got.Fields {
			fmt.Println(f.At(i))
		}
	}

}
