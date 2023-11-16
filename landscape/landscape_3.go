package landscape

type Landscape3 []uint

// TuringCompleteExample3 returns a Landscape with a rising and falling peak
func TuringCompleteExample3() Landscape3 {
	return Landscape3{4, 6, 1, 4, 6, 5, 1, 4, 1, 2, 6, 5, 6, 1, 4, 2}
	//               0, 0, 5, 2, 0, 1, 5, 2, 5, 4, 0, 1, 0, 3, 0, 0   = 28 water volume
}

// StaggeredDecliningPeaks3 returns a Landscape with a rising peak and then staggered declining peaks
func StaggeredDecliningPeaks3() Landscape3 {
	return Landscape3{4, 6, 1, 4, 6, 5, 1, 4, 1, 2, 6, 4, 5, 1, 4, 2}
	//               0, 0, 5, 2, 0, 1, 5, 2, 5, 4, 0, 1, 0, 3, 0, 0   = 28 water volume
}

func (l Landscape3) Water() (totalWater uint) {
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
			totalWater += leftMax - l[i]
		}
	}

	// Calculate water trapped in the right section
	var rightMax uint
	for i := len(l) - 1; i > highestIdx; i-- {
		if l[i] > rightMax {
			rightMax = l[i]
		} else {
			totalWater += rightMax - l[i]
		}
	}

	return totalWater
}
