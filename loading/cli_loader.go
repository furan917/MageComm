package loading

import (
	"fmt"
	"time"
)

func Indicator(stopLoading chan bool) {
	loadingChars := []string{"-", "\\", "|", "/"}
	i := 0

	for {
		select {
		case <-stopLoading:
			fmt.Print("\r")
			return
		default:
			fmt.Printf("\rWaiting on response... %s", loadingChars[i])
			i = (i + 1) % len(loadingChars)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
