package landscape

type Landscape2[N Number] []N

// TuringCompleteExample2 returns a Landscape with a rising and falling peak
func TuringCompleteExample2[N Number]() Landscape2[N] {
	return Landscape2[N]{4, 6, 1, 4, 6, 5, 1, 4, 1, 2, 6, 5, 6, 1, 4, 2}
	//               0, 0, 5, 2, 0, 1, 5, 2, 5, 4, 0, 1, 0, 3, 0, 0   = 28 water volume
}

// StaggeredDecliningPeaks2 returns a Landscape with a rising peak and then staggered declining peaks
func StaggeredDecliningPeaks2[N Number]() Landscape2[N] {
	return Landscape2[N]{4, 6, 1, 4, 6, 5, 1, 4, 1, 2, 6, 4, 5, 1, 4, 2}
	//               0, 0, 5, 2, 0, 1, 5, 2, 5, 4, 0, 1, 0, 3, 0, 0   = 28 water volume
}

// Water returns the total amount of water that could be held in the landscape
func (l Landscape2[N]) Water() (water N) {
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
