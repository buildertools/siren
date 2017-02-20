package demo

import (
	"github.com/buildertools/siren"
	"context"
	"math"
	"math/rand"
	"time"
)

func init() {
	siren.RegisterBackend(`randoRate`, RandoRate20)
	siren.RegisterBackend(`randoCounter`, RandoCounter20)
}

var rateSeries []float64
var durationSeries []float64
var counterSeries []float64

func Start() {
	if rateSeries == nil {
		rateSeries = []float64{}
		go func() {
			t := time.NewTicker(1 * time.Second)
			for {
				select {
				case <-t.C:
					if len(rateSeries) == 20 {
						rateSeries = rateSeries[1:len(rateSeries)]
					}
					rateSeries = append(rateSeries, rand.Float64())
				}
			}
		}()
	}
	if counterSeries == nil {
		counterSeries = []float64{}
		go func() {
			t := time.NewTicker(1 * time.Second)
			for {
				select {
				case <-t.C:
					if len(counterSeries) == 20 {
						counterSeries = counterSeries[1:len(counterSeries)]
					}
					counterSeries = append(counterSeries, math.Trunc(rand.Float64()))
				}
			}
		}()
	}

}

func r(in []float64) []float64 {
	o := make([]float64, len(in))
	for i, j := 0, len(in) - 1; i < j; i, j = i + 1, j - 1 {
		o[i], o[j] = in[j], in[i]
	}
	return o
}

func RandoCounter20(c context.Context, q string) ([]float64, error) {
	return r(counterSeries), nil
}

func RandoRate20(c context.Context, q string) ([]float64, error) {
	return r(rateSeries), nil
}

func RandoDuration20(c context.Context, q string) ([]float64, error) {
	return []float64{
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
		rand.NormFloat64(),
	}, nil
}
