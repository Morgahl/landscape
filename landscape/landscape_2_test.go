package landscape_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/morgahl/play/landscape"
)

const (
	MAX_TEST_ARRAY_SIZE = 1_000_000 // TODO: make this a flag

	RAND_SEED = 0 // TODO: make this a flag
)

func BenchmarkLandscape2_TuringCompleteExample2(b *testing.B) {
	b.Helper()
	benchLand2Cons_N_T_Validate(b, landscape.TuringCompleteExample2[int](), 28, "_16")
}

func BenchmarkLandscape2_StaggeredDecliningPeaks2(b *testing.B) {
	b.Helper()
	benchLand2Cons_N_T_Validate(b, landscape.StaggeredDecliningPeaks2[int](), 28, "_16")
}

func benchLand2Cons_N_T_Validate[N landscape.Number](b *testing.B, l landscape.Landscape2[N], expect N, suffix string) {
	b.Helper()
	b.Run(fmt.Sprintf("[%T]%s", N(0), suffix), benchmarkLandscape2_N_Validate(l, expect))
	b.Run(fmt.Sprintf("[%T]%s_perm", N(0), suffix), benchmarkLandscape2_N(permute2(l)))
}

func benchmarkLandscape2_N_Validate[N landscape.Number](l landscape.Landscape2[N], expect N) func(b *testing.B) {
	return func(b *testing.B) {
		var water N
		for i := 0; i < b.N; i++ {
			water = l.Water()
		}
		b.StopTimer()
		if water != N(expect) {
			b.StartTimer()
			b.Fatalf("Expected %v, got %v", expect, water)
		}
	}
}

func BenchmarkLandscape2_Incline2(b *testing.B) {
	b.Helper()
	for n := 10; n <= MAX_TEST_ARRAY_SIZE; n *= 10 {
		benchLand2Cons_N_T(b, generateIncline2[int](n, 1), fmt.Sprintf("_%d_0", n))
		benchLand2Cons_N_T(b, generateIncline2[int](n, 2), fmt.Sprintf("_%d_1", n))
	}
}

func BenchmarkLandscape2_Decline2(b *testing.B) {
	b.Helper()
	for n := 10; n <= MAX_TEST_ARRAY_SIZE; n *= 10 {
		benchLand2Cons_N_T(b, generateDecline2[int](n, 1), fmt.Sprintf("_%d_0", n))
		benchLand2Cons_N_T(b, generateDecline2[int](n, 2), fmt.Sprintf("_%d_1", n))
	}
}

func BenchmarkLandscape2_InclineDecline2(b *testing.B) {
	b.Helper()
	for n := 10; n <= MAX_TEST_ARRAY_SIZE; n *= 10 {
		benchLand2Cons_N_T(b, generateInclineDecline2[int](n, 1), fmt.Sprintf("_%d_0", n))
		benchLand2Cons_N_T(b, generateInclineDecline2[int](n, 2), fmt.Sprintf("_%d_1", n))
	}
}

func BenchmarkLandscape2_DeclineIncline2(b *testing.B) {
	b.Helper()
	for n := 10; n <= MAX_TEST_ARRAY_SIZE; n *= 10 {
		benchLand2Cons_N_T(b, generateDeclineIncline2[int](n, 1), fmt.Sprintf("_%d_0", n))
		benchLand2Cons_N_T(b, generateDeclineIncline2[int](n, 2), fmt.Sprintf("_%d_1", n))
	}
}

func BenchmarkLandscape2_Cross2(b *testing.B) {
	b.Helper()
	for n := 10; n <= MAX_TEST_ARRAY_SIZE; n *= 10 {
		benchLand2Cons_N_T(b, generateCross2[int](n), fmt.Sprintf("_%d", n))
	}
}

func benchLand2Cons_N_T[N landscape.Number](b *testing.B, l landscape.Landscape2[N], suffix string) {
	b.Helper()
	b.Run(fmt.Sprintf("[%T]%s", N(0), suffix), benchmarkLandscape2_N(l))
	b.Run(fmt.Sprintf("[%T]%s_perm", N(0), suffix), benchmarkLandscape2_N(permute2(l)))
}

func benchmarkLandscape2_N[N landscape.Number](l landscape.Landscape2[N]) func(b *testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l.Water()
		}
	}
}

func generateIncline2[N landscape.Number](n, step int) landscape.Landscape2[N] {
	if step < 1 {
		step = 1
	}
	l := make(landscape.Landscape2[N], n)
	for i, j := 0, 0; i < n; i, j = i+step, j+1 {
		l[i] = convAbs[N](j)
	}
	return l
}

func generateDecline2[N landscape.Number](n, step int) landscape.Landscape2[N] {
	if step < 1 {
		step = 1
	}
	l := make(landscape.Landscape2[N], n)
	for i, j := 0, n/step; i < n; i, j = i+step, j-1 {
		l[i] = convAbs[N](j)
	}
	return l
}

func generateInclineDecline2[N landscape.Number](n, step int) landscape.Landscape2[N] {
	if step < 1 {
		step = 1
	}
	mid := n / 2
	l := make(landscape.Landscape2[N], n)
	for i, j := 0, 0; i < mid; i, j = i+step, j+1 {
		l[i] = convAbs[N](j)
	}
	for i, j := mid, n/step; i < n; i, j = i+step, j-1 {
		l[i] = convAbs[N](j)
	}
	return l
}

func generateDeclineIncline2[N landscape.Number](n, step int) landscape.Landscape2[N] {
	if step < 1 {
		step = 1
	}
	mid := n / 2
	l := make(landscape.Landscape2[N], n)
	for i, j := 0, n/step; i < mid; i, j = i+step, j-1 {
		l[i] = convAbs[N](j)
	}
	for i, j := mid, 0; i < n; i, j = i+step, j+1 {
		l[i] = convAbs[N](j)
	}
	return l
}

func generateCross2[N landscape.Number](n int) landscape.Landscape2[N] {
	l := make(landscape.Landscape2[N], n)
	mid := n / 2
	for i, j, k := 0, 0, mid; i < n; i, j, k = i+2, j+1, k-1 {
		l[i] = convAbs[N](j)
		l[i+1] = convAbs[N](k)
	}
	return l
}

func permute2[N landscape.Number](l landscape.Landscape2[N]) landscape.Landscape2[N] {
	length := len(l)
	lN := make(landscape.Landscape2[N], length)
	copy(lN, l)
	rand := rand.New(rand.NewSource(RAND_SEED))
	for i := range lN {
		j := rand.Intn(length - i)
		lN[i], lN[i+j] = lN[i+j], lN[i]
	}
	return lN
}

func convAbs[N landscape.Number](i int) N {
	n := N(i)
	if n < 0 {
		return -n
	}
	return n
}
