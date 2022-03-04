package main

import (
	"time"

	"github.com/marcusolsson/grafana-ynab-datasource/pkg/ynab"
)

type Range struct {
	Start time.Time
	End   time.Time
}

func (r Range) Contains(t time.Time) bool {
	return r.Start.UnixNano() <= t.UnixNano() && t.UnixNano() < r.End.UnixNano()
}

type Bucket struct {
	Range        Range
	Measurements []Measurement
}

type Histogram struct {
	buckets []Bucket
	series  TimeSeries
}

func NewHistogram(series TimeSeries) *Histogram {
	return &Histogram{
		buckets: make([]Bucket, 0),
		series:  series,
	}
}

type Period int

const (
	PeriodDaily Period = iota
	PeriodWeekly
	PeriodMonthly
)

func (g *Histogram) Fill(start, end time.Time, p Period) {
	switch p {
	case PeriodDaily:
		first := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
		last := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)

		for b := first; b != last; b = b.AddDate(0, 0, 1) {
			g.buckets = append(g.buckets, Bucket{
				Range: Range{
					Start: b,
					End:   b.AddDate(0, 0, 1),
				},
				Measurements: []Measurement{},
			})
		}
	case PeriodWeekly:
		first := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
		for {
			if first.Weekday() == time.Monday {
				break
			}
			first = first.AddDate(0, 0, -1)
		}

		last := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
		for {
			if last.Weekday() == time.Monday {
				break
			}
			last = last.AddDate(0, 0, 1)
		}

		for b := first; b != last; b = b.AddDate(0, 0, 7) {
			g.buckets = append(g.buckets, Bucket{
				Range: Range{
					Start: b,
					End:   b.AddDate(0, 0, 7),
				},
				Measurements: []Measurement{},
			})
		}
	case PeriodMonthly:
		first := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
		last := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)

		for b := first; b != last; b = b.AddDate(0, 1, 0) {
			g.buckets = append(g.buckets, Bucket{
				Range: Range{
					Start: b,
					End:   b.AddDate(0, 1, 0),
				},
				Measurements: []Measurement{},
			})
		}
	}
}

func (g *Histogram) Add(m Measurement) {
	for j, bucket := range g.buckets {
		if bucket.Range.Contains(m.Time) {
			g.buckets[j].Measurements = append(g.buckets[j].Measurements, m)
		}
	}
}

func (g *Histogram) EachBucket(reducer func(r Range, ms []Measurement) Measurement, fn func(m Measurement), gapFill string) {
	var lastMeasurement Measurement
	for i, b := range g.buckets {
		if len(b.Measurements) > 0 {
			lastMeasurement = reducer(b.Range, b.Measurements)
			fn(lastMeasurement)
		} else {
			if i > 0 {
				if gapFill == "last" {
					fn(Measurement{
						Time:   b.Range.Start,
						Value:  lastMeasurement.Value,
						Labels: lastMeasurement.Labels,
					})
				}
			} else {
				if gapFill == "last" {
					fn(Measurement{
						Time:   b.Range.Start,
						Value:  0,
						Labels: lastMeasurement.Labels,
					})
				}
			}
		}
	}
}

type Measurement struct {
	Time   time.Time
	Value  float64
	Labels map[string]string
}

type TimeSeries interface {
	Time(i int) time.Time
	Value(i int) float64
	Len() int
	Labels(i int) map[string]string
}

type TimeSeriesTransactions []ynab.Transaction

func (a TimeSeriesTransactions) Len() int {
	return len(a)
}

func (a TimeSeriesTransactions) Time(i int) time.Time {
	t, err := time.Parse("2006-01-02", a[i].Date)
	if err != nil {
		panic(err)
	}
	return t
}

func (a TimeSeriesTransactions) Value(i int) float64 {
	return float64(a[i].Amount)
}

func (a TimeSeriesTransactions) Labels(i int) map[string]string {
	return map[string]string{
		"account_id":    a[i].AccountID,
		"account_name":  a[i].AccountName,
		"payee_id":      a[i].PayeeID,
		"payee_name":    a[i].PayeeName,
		"category_id":   a[i].CategoryID,
		"category_name": a[i].CategoryName,
	}
}

var _ TimeSeries = TimeSeriesTransactions{}

type TimeSeriesBalance []ynab.Balance

func (a TimeSeriesBalance) Len() int {
	return len(a)
}

func (a TimeSeriesBalance) Time(i int) time.Time {
	t, err := time.Parse("2006-01-02", a[i].Date)
	if err != nil {
		panic(err)
	}
	return t
}

func (a TimeSeriesBalance) Value(i int) float64 {
	return float64(a[i].Amount)
}

func (a TimeSeriesBalance) Labels(i int) map[string]string {
	return map[string]string{
		"account_id":   a[i].AccountID,
		"account_name": a[i].AccountName,
	}
}

var _ TimeSeries = TimeSeriesBalance{}

func Regularize(series TimeSeries, p Period, aligner func(r Range, ms []Measurement) Measurement, gapFill string) ([]Measurement, error) {
	if series.Len() < 0 {
		return []Measurement{}, nil
	}

	hist := NewHistogram(series)
	hist.Fill(series.Time(0), series.Time(series.Len()-1), p)

	for i := 0; i < series.Len(); i++ {
		m := Measurement{
			Time:   series.Time(i),
			Value:  series.Value(i),
			Labels: series.Labels(i),
		}
		hist.Add(m)
	}

	var res []Measurement

	hist.EachBucket(
		aligner,
		func(m Measurement) {
			res = append(res, m)
		},
		gapFill,
	)

	return res, nil
}

func alignLast(r Range, ms []Measurement) Measurement {
	last := ms[len(ms)-1]
	last.Time = r.Start
	return last
}

func alignTotal(r Range, ms []Measurement) Measurement {
	var total float64

	for _, m := range ms {
		total += m.Value
	}

	return Measurement{
		Time:   r.Start,
		Value:  total,
		Labels: ms[0].Labels,
	}
}
