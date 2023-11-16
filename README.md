# A random journey through Go and a simple problem

I wanted to write a simple solution to the following code challenge but ended up finding my way down an Github Copilot and ChatGpt rabbit hole....

Let's write a Go program that receives a slice of uint integers that represent landscape elevations. I want to performantly calculate the total amount of water units that can be contained by this landscape. Assume that water will fill evenly to each peak and that no water can be contained by beginning and ending bounds.

The example landscape can be visualized as follows:

```plain
[4, 6, 1, 4, 6, 5, 1, 4, 1, 2, 6, 5, 6, 1, 4, 2]
 0, 0, 5, 2, 0, 1, 5, 2, 5, 4, 0, 1, 0, 3, 0, 0 == 28 water volume
```

Or my Ascii version....

```plain
| | | | | | | | | | | | | | | | |
| | | | | | | | | | | | | | | | |
| |#| | |#| | | | | |#| |#| | | |
| |#| | |#|#| | | | |#|#|#| | | |
|#|#| |#|#|#| |#| | |#|#|#| |#| |
|#|#| |#|#|#| |#| | |#|#|#| |#| |
|#|#| |#|#|#| |#| |#|#|#|#| |#|#|
|#|#|#|#|#|#|#|#|#|#|#|#|#|#|#|#|
=================================
|4|6|1|4|6|5|1|4|1|2|6|5|6|1|4|2|
```

and with a little bit of imagination....

```plain
| | | | | | | | | | | | | | | | |
| | | | | | | | | | | | | | | | |
| |#|~|~|#|~|~|~|~|~|#|~|#| | | |
| |#|~|~|#|#|~|~|~|~|#|#|#| | | |
|#|#|~|#|#|#|~|#|~|~|#|#|#|~|#| |
|#|#|~|#|#|#|~|#|~|~|#|#|#|~|#| |
|#|#|~|#|#|#|~|#|~|#|#|#|#|~|#|#|
|#|#|#|#|#|#|#|#|#|#|#|#|#|#|#|#|
=================================
|0|0|5|2|0|1|5|2|5|4|0|1|0|3|0|0|
```


My initial attempt is as follows:

```go
type Landscape []uint

// Water returns the total amount of water that could be held in the landscape
func (l Landscape) Water() (water uint) {
	for startIdx := 0; ; {
		switch peak := l.nextPeak(startIdx); peak {
		default:
			water += l.waterBetween(startIdx, peak)
			startIdx = peak

		case -1:
			return water
		}
	}
}

type peakWalkerState uint

const (
	LEVEL_OR_INCLINING peakWalkerState = iota
	HAS_DECLINED
	HAS_VALLEY
)

// nextPeak returns the index of the next peak in the landscape
// TODO: this has (O(n^2)) for declining staggered peaks, maybe we iterate from both ends of the landscape?
func (l Landscape) nextPeak(startIdx int) (peakIdx int) {
	var walkerState peakWalkerState = LEVEL_OR_INCLINING
	var waterIdx int // the index of the last peak not exceeding the startIdx peak when walkerState == HAS_VALLEY
	startHeight := l[startIdx]
	for peakIdx = startIdx + 1; peakIdx < len(l); peakIdx++ {
		peakValue := l[peakIdx]
		if startHeight <= peakValue {
			return peakIdx
		}

		if walkerState == HAS_VALLEY && l[waterIdx] <= peakValue {
			waterIdx = peakIdx
			continue
		}

		lastValue := l[peakIdx-1]
		if walkerState == HAS_DECLINED && lastValue < peakValue {
			walkerState = HAS_VALLEY
			waterIdx = peakIdx
			continue
		}

		if walkerState == LEVEL_OR_INCLINING && peakValue < lastValue {
			walkerState = HAS_DECLINED
			continue
		}
	}

	if walkerState == HAS_VALLEY {
		return waterIdx
	}

	return -1
}

// waterBetween returns the amount of water that could be held between two peaks
func (l Landscape) waterBetween(startIdx, endIdx int) uint {
	var water uint
	maxHeight := l[startIdx]
	if l[endIdx] < maxHeight {
		maxHeight = l[endIdx]
	}

	for i := startIdx + 1; i < endIdx; i++ {
		if l[i] < maxHeight {
			water += maxHeight - l[i]
		}
	}

	return water
}
```

I next asked Github Copilot to help me write a more efficient version of this code. I gave it the following prompt:

```
// Using Landscape as an example create a variant that iterates inwards from both ends of the
// landscape and returns the total amount of water that could be held in the landscape (O(n)) time
// and O(1) space complexity. Implemented this as Landscape2 to avoid breaking the existing code.

```

It came up with the following:

```go
type Landscape2 []uint

// Water returns the total amount of water that could be held in the landscape
func (l Landscape2) Water() (water uint) {
	left := 0
	right := len(l) - 1
	leftMax := l[left]
	rightMax := l[right]

	for left < right {
		if l[left] < l[right] {
			if l[left] > leftMax {
				leftMax = l[left]
			} else {
				water += leftMax - l[left]
			}
			left++
		} else {
			if l[right] > rightMax {
				rightMax = l[right]
			} else {
				water += rightMax - l[right]
			}
			right--
		}
	}

	return water
}
```

with the following results with this command `go test --bench . ./... --count 3 --benchmem`:

```sh
goos: linux
goarch: amd64
pkg: github.com/morgahl/play/landscape
cpu: Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz
BenchmarkTuringCompleteExample2-16              68239588                17.74 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2-16              68607103                17.72 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2-16              68229889                17.68 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2-16            71797387                16.63 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2-16            72100652                16.72 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2-16            72960302                16.64 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample-16               30889143                36.96 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample-16               32124814                37.02 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample-16               32204390                36.89 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks-16             26364007                44.99 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks-16             26525901                45.03 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks-16             26318559                45.38 ns/op            0 B/op          0 allocs/op
PASS`
ok      github.com/morgahl/play/landscape       14.715s
```

As you can see the new version is significantly faster than the original. The original version is O(n^2) and the new version is O(n).


Next I asked Chat GPT to write a function that would return the total amount of water that could be held in the landscape. I gave it the following prompt:

```
I am writing a Go program to calculate the total amount of water that could be held in a landscape. The landscape is represented by a slice of uint values. The value of each element in the slice represents the height of the landscape at that point. The landscape is bounded by the edges of the slice. The edges of the slice do not retain water. The function should return the total amount of water units that could be held in the landscape.
```

It came up with the following which since I cannot keep the ChatGPT history I will copy here what I transposed into the same effective interface as the original:

```go
type Landscape3 []uint

func (l Landscape3) Water() uint {
	var totalWater uint = 0

	// Find the highest point in the landscape
	var highestIdx int
	for i, height := range l {
		if height > l[highestIdx] {
			highestIdx = i
		}
	}

	// Calculate water trapped in the left section
	var leftMax uint
	for i := 0; i < highestIdx; i++ {
		if l[i] > leftMax {
			leftMax = l[i]
		} else {
			totalWater += uint(leftMax - l[i])
		}
	}

	// Calculate water trapped in the right section
	var rightMax uint
	for i := len(l) - 1; i > highestIdx; i-- {
		if l[i] > rightMax {
			rightMax = l[i]
		} else {
			totalWater += uint(rightMax - l[i])
		}
	}

	return totalWater
}
```

This resulted in the following benchmark results:

```sh
goos: linux
goarch: amd64
pkg: github.com/morgahl/play/landscape
cpu: Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz
BenchmarkTuringCompleteExample2-16                      67432398                17.03 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2-16                      75874927                17.80 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2-16                      66934404                17.71 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_10-16                    8644672               138.2 ns/op             0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_10-16                    8556210               138.9 ns/op             0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_10-16                    8645508               138.4 ns/op             0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_100-16                    877584              1322 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_100-16                    882547              1320 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_100-16                    868959              1335 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_1_000-16                   90420             13123 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_1_000-16                   90741             13172 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_1_000-16                   91197             13148 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_10_000-16                   8852            134462 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_10_000-16                   8854            134802 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_10_000-16                   8880            136476 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_100_000-16                   904           1334029 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_100_000-16                   889           1329172 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_100_000-16                   904           1328104 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_1_000_000-16                  86          13356257 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_1_000_000-16                  86          13749585 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2_1_000_000-16                  86          13507981 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2-16                    71707718                16.55 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2-16                    67601444                16.50 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2-16                    71848972                16.59 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_10-16                  8193094               140.9 ns/op             0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_10-16                  8384701               139.4 ns/op             0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_10-16                  8314308               142.0 ns/op             0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_100-16                  874182              1336 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_100-16                  879804              1325 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_100-16                  882579              1323 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_1_000-16                 88669             13431 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_1_000-16                 88418             13135 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_1_000-16                 89780             13279 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_10_000-16                 8823            131315 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_10_000-16                 8839            131176 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_10_000-16                 8850            131378 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_100_000-16                 908           1336785 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_100_000-16                 909           1317692 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_100_000-16                 904           1319856 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_1_000_000-16                87          13342330 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_1_000_000-16                90          13310906 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2_1_000_000-16                86          13306466 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3-16                      46443402                25.69 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3-16                      46414300                25.73 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3-16                      48485239                25.73 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_10-16                    5698378               209.5 ns/op             0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_10-16                    5598609               209.6 ns/op             0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_10-16                    5623002               211.0 ns/op             0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_100-16                    573964              1986 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_100-16                    548678              1984 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_100-16                    591237              1992 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_1_000-16                   59584             19671 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_1_000-16                   60763             19727 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_1_000-16                   60273             20081 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_10_000-16                   5822            196678 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_10_000-16                   5964            197182 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_10_000-16                   5704            198154 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_100_000-16                   604           1969751 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_100_000-16                   607           1973796 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_100_000-16                   606           1978030 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_1_000_000-16                  58          19967522 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_1_000_000-16                  57          19933718 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3_1_000_000-16                  57          20269061 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3-16                    42855152                25.77 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3-16                    46292902                26.23 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3-16                    46284331                25.81 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_10-16                  5741912               209.2 ns/op             0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_10-16                  5737845               208.7 ns/op             0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_10-16                  5720618               208.7 ns/op             0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_100-16                  590768              1990 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_100-16                  593277              1986 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_100-16                  586759              1994 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_1_000-16                 60768             19683 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_1_000-16                 60712             19694 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_1_000-16                 59349             19975 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_10_000-16                 5978            197400 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_10_000-16                 5911            197206 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_10_000-16                 5912            197207 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_100_000-16                 592           1969020 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_100_000-16                 607           2003030 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_100_000-16                 590           1975902 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_1_000_000-16                57          20055581 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_1_000_000-16                55          20222576 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3_1_000_000-16                57          19916983 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample-16                       31500775                37.86 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample-16                       31419698                37.74 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample-16                       30790783                38.69 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_10-16                     3585272               333.0 ns/op             0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_10-16                     3377476               329.8 ns/op             0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_10-16                     3615742               331.7 ns/op             0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_100-16                     365022              3271 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_100-16                     362244              3283 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_100-16                     363118              3273 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_1_000-16                    36108             33142 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_1_000-16                    35599             32526 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_1_000-16                    36777             32653 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_10_000-16                    3628            326165 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_10_000-16                    3622            327123 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_10_000-16                    3626            334235 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_100_000-16                    352           3326712 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_100_000-16                    358           3364580 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_100_000-16                    364           3298941 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_1_000_000-16                   31          32925020 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_1_000_000-16                   33          32684537 ns/op               0 B/op          0 allocs/op
BenchmarkTuringCompleteExample_1_000_000-16                   36          32739239 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks-16                     26358043                46.43 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks-16                     25431595                45.53 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks-16                     26315384                45.67 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_10-16                   3180616               367.1 ns/op             0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_10-16                   3264748               382.9 ns/op             0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_10-16                   3249458               357.0 ns/op             0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_100-16                   347844              3216 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_100-16                   350271              3404 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_100-16                   349630              3423 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_1_000-16                  34604             34473 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_1_000-16                  35188             34545 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_1_000-16                  33661             34390 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_10_000-16                  3506            345978 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_10_000-16                  3513            344671 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_10_000-16                  3447            346794 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_100_000-16                  344           3462802 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_100_000-16                  339           3458596 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_100_000-16                  345           3453944 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_1_000_000-16                 33          34589158 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_1_000_000-16                 33          34536442 ns/op               0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks_1_000_000-16                 33          34665264 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/morgahl/play/landscape       163.843s
```

For a lark I decided to try out what it would look like if these landscapes were declared as:

```go
type Landscape [16]uint
```

Which gave the following results:

```sh
goos: linux
goarch: amd64
pkg: github.com/morgahl/play/landscape
cpu: Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz
BenchmarkTuringCompleteExample2-16              82682785                13.31 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2-16              92530477                13.38 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample2-16              92944720                13.33 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2-16            81514539                14.91 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2-16            81496824                13.03 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks2-16            80355169                14.77 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3-16              39812613                29.16 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3-16              42176443                28.90 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample3-16              39721681                28.98 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3-16            41321030                27.90 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3-16            43306758                29.46 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks3-16            40444617                29.39 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample-16               32459098                36.30 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample-16               32411492                36.24 ns/op            0 B/op          0 allocs/op
BenchmarkTuringCompleteExample-16               32318614                36.28 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks-16             26822952                44.31 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks-16             26808928                44.31 ns/op            0 B/op          0 allocs/op
BenchmarkStaggeredDecliningPeaks-16             26924964                46.03 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/morgahl/play/landscape       21.870s
```

## Conclusion

My initial attempt was well intentioned but misguided in several ways. I was trying to be clever and only interate from one side but this leads to cascading lookups causing O(N^2) worst case scenarios for declinging staggered peaks.

The Github Copilot suggestion was a good one and was more performant then my approach due to all calculations being done in a single pass.

The ChatGPT suggestion was also a good one and was more performant then my approach but less then the Github Copilot suggestion. This is likely due to the initial pass to find the highest (pivot) point.