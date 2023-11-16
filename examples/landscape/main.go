package main

import (
	"fmt"

	"github.com/morgahl/play/landscape"
)

func main() {
	tce := landscape.TuringCompleteExample()
	fmt.Printf("TuringCompleteExample: %v\n", tce)
	fmt.Printf("Water: %d\n", tce.Water())

	sdp := landscape.StaggeredDecliningPeaks()
	fmt.Printf("StaggeredDecliningPeaks: %v\n", sdp)
	fmt.Printf("Water: %d\n", sdp.Water())

	tce2 := landscape.TuringCompleteExample2[int]()
	fmt.Printf("TuringCompleteExample2: %v\n", tce2)
	fmt.Printf("Water: %d\n", tce2.Water())

	sdp2 := landscape.StaggeredDecliningPeaks2[int]()
	fmt.Printf("StaggeredDecliningPeaks2: %v\n", sdp2)
	fmt.Printf("Water: %d\n", sdp2.Water())

	tce3 := landscape.TuringCompleteExample3()
	fmt.Printf("TuringCompleteExample3: %v\n", tce3)
	fmt.Printf("Water: %d\n", tce3.Water())

	sdp3 := landscape.StaggeredDecliningPeaks3()
	fmt.Printf("StaggeredDecliningPeaks3: %v\n", sdp3)
	fmt.Printf("Water: %d\n", sdp3.Water())
}
