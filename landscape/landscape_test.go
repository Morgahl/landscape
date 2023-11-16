package landscape_test

import (
	"testing"

	"github.com/morgahl/play/landscape"
)

func BenchmarkTuringCompleteExample_1(b *testing.B) {
	tce := landscape.TuringCompleteExample()
	for i := 0; i < b.N; i++ {
		tce.Water()
	}
}

func BenchmarkTuringCompleteExample_10(b *testing.B)     { benchmarkTuringCompleteExample_N(b, 10) }
func BenchmarkTuringCompleteExample_100(b *testing.B)    { benchmarkTuringCompleteExample_N(b, 100) }
func BenchmarkTuringCompleteExample_1_000(b *testing.B)  { benchmarkTuringCompleteExample_N(b, 1_000) }
func BenchmarkTuringCompleteExample_10_000(b *testing.B) { benchmarkTuringCompleteExample_N(b, 10_000) }
func BenchmarkTuringCompleteExample_100_000(b *testing.B) {
	benchmarkTuringCompleteExample_N(b, 100_000)
}
func BenchmarkTuringCompleteExample_1_000_000(b *testing.B) {
	benchmarkTuringCompleteExample_N(b, 1_000_000)
}

func benchmarkTuringCompleteExample_N(b *testing.B, n int) {
	b.StopTimer()
	tce := landscape.TuringCompleteExample()
	tce10k := make(landscape.Landscape, n*len(tce))
	for i := 0; i < n; i++ {
		offset := i * len(tce)
		copy(tce10k[offset:offset+len(tce)], tce)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tce10k.Water()
	}
}

func BenchmarkStaggeredDecliningPeaks_1(b *testing.B) {
	sdc := landscape.StaggeredDecliningPeaks()
	for i := 0; i < b.N; i++ {
		sdc.Water()
	}
}

func BenchmarkStaggeredDecliningPeaks_10(b *testing.B)  { benchmarkStaggeredDecliningPeaks_N(b, 10) }
func BenchmarkStaggeredDecliningPeaks_100(b *testing.B) { benchmarkStaggeredDecliningPeaks_N(b, 100) }
func BenchmarkStaggeredDecliningPeaks_1_000(b *testing.B) {
	benchmarkStaggeredDecliningPeaks_N(b, 1_000)
}
func BenchmarkStaggeredDecliningPeaks_10_000(b *testing.B) {
	benchmarkStaggeredDecliningPeaks_N(b, 10_000)
}
func BenchmarkStaggeredDecliningPeaks_100_000(b *testing.B) {
	benchmarkStaggeredDecliningPeaks_N(b, 100_000)
}
func BenchmarkStaggeredDecliningPeaks_1_000_000(b *testing.B) {
	benchmarkStaggeredDecliningPeaks_N(b, 1_000_000)
}

func benchmarkStaggeredDecliningPeaks_N(b *testing.B, n int) {
	b.StopTimer()
	sdc := landscape.StaggeredDecliningPeaks()
	sdc10k := make(landscape.Landscape, n*len(sdc))
	for i := 0; i < n; i++ {
		offset := i * len(sdc)
		copy(sdc10k[offset:offset+len(sdc)], sdc)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sdc10k.Water()
	}
}
