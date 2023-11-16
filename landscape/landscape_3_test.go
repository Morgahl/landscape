package landscape_test

import (
	"testing"

	"github.com/morgahl/play/landscape"
)

func BenchmarkTuringCompleteExample3_1(b *testing.B) {
	tce := landscape.TuringCompleteExample3()
	for i := 0; i < b.N; i++ {
		tce.Water()
	}
}

func BenchmarkTuringCompleteExample3_10(b *testing.B)    { benchmarkTuringCompleteExample3_N(b, 10) }
func BenchmarkTuringCompleteExample3_100(b *testing.B)   { benchmarkTuringCompleteExample3_N(b, 100) }
func BenchmarkTuringCompleteExample3_1_000(b *testing.B) { benchmarkTuringCompleteExample3_N(b, 1_000) }
func BenchmarkTuringCompleteExample3_10_000(b *testing.B) {
	benchmarkTuringCompleteExample3_N(b, 10_000)
}
func BenchmarkTuringCompleteExample3_100_000(b *testing.B) {
	benchmarkTuringCompleteExample3_N(b, 100_000)
}
func BenchmarkTuringCompleteExample3_1_000_000(b *testing.B) {
	benchmarkTuringCompleteExample3_N(b, 1_000_000)
}

func benchmarkTuringCompleteExample3_N(b *testing.B, n int) {
	b.StopTimer()
	tce := landscape.TuringCompleteExample3()
	tce10k := make(landscape.Landscape3, n*len(tce))
	for i := 0; i < n; i++ {
		offset := i * len(tce)
		copy(tce10k[offset:offset+len(tce)], tce)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tce10k.Water()
	}
}

func BenchmarkStaggeredDecliningPeaks3_1(b *testing.B) {
	sdc := landscape.StaggeredDecliningPeaks3()
	for i := 0; i < b.N; i++ {
		sdc.Water()
	}
}

func BenchmarkStaggeredDecliningPeaks3_10(b *testing.B)  { benchmarkStaggeredDecliningPeaks3_N(b, 10) }
func BenchmarkStaggeredDecliningPeaks3_100(b *testing.B) { benchmarkStaggeredDecliningPeaks3_N(b, 100) }
func BenchmarkStaggeredDecliningPeaks3_1_000(b *testing.B) {
	benchmarkStaggeredDecliningPeaks3_N(b, 1_000)
}
func BenchmarkStaggeredDecliningPeaks3_10_000(b *testing.B) {
	benchmarkStaggeredDecliningPeaks3_N(b, 10_000)
}
func BenchmarkStaggeredDecliningPeaks3_100_000(b *testing.B) {
	benchmarkStaggeredDecliningPeaks3_N(b, 100_000)
}
func BenchmarkStaggeredDecliningPeaks3_1_000_000(b *testing.B) {
	benchmarkStaggeredDecliningPeaks3_N(b, 1_000_000)
}

func benchmarkStaggeredDecliningPeaks3_N(b *testing.B, n int) {
	b.StopTimer()
	sdc := landscape.StaggeredDecliningPeaks3()
	sdc10k := make(landscape.Landscape3, n*len(sdc))
	for i := 0; i < n; i++ {
		offset := i * len(sdc)
		copy(sdc10k[offset:offset+len(sdc)], sdc)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sdc10k.Water()
	}
}
