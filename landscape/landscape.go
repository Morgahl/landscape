package landscape

type Landscape []uint

// TuringCompleteExample returns a Landscape with a rising and falling peak
func TuringCompleteExample() Landscape {
	return Landscape{4, 6, 1, 4, 6, 5, 1, 4, 1, 2, 6, 5, 6, 1, 4, 2}
	//               0, 0, 5, 2, 0, 1, 5, 2, 5, 4, 0, 1, 0, 3, 0, 0   = 28 water volume
	//               0  1        4                10    12    14   -1 : 6 peak segments
}

// StaggeredDecliningPeaks returns a Landscape with a rising peak and then staggered declining peaks
func StaggeredDecliningPeaks() Landscape {
	return Landscape{4, 6, 1, 4, 6, 5, 1, 4, 1, 2, 6, 4, 5, 1, 4, 2}
	//               0, 0, 5, 2, 0, 1, 5, 2, 5, 4, 0, 1, 0, 3, 0, 0   = 28 water volume
	//               0  1        4                10    12    14   -1 : 6 peak segments
}

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
	for peakIdx = startIdx + 1; peakIdx < len(l); peakIdx++ {
		if l[startIdx] <= l[peakIdx] {
			return peakIdx
		}

		if walkerState == HAS_VALLEY && l[waterIdx] <= l[peakIdx] {
			waterIdx = peakIdx
			continue
		}

		if walkerState == HAS_DECLINED && l[peakIdx-1] < l[peakIdx] {
			walkerState = HAS_VALLEY
			waterIdx = peakIdx
			continue
		}

		if walkerState == LEVEL_OR_INCLINING && l[peakIdx] < l[peakIdx-1] {
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
func (l Landscape) waterBetween(startIdx, endIdx int) (water uint) {
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
