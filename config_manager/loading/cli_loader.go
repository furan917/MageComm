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
			fmt.Print("\r\x1b[K") // clear line
			return
		default:
			fmt.Printf("\rWaiting on response: Press ENTER to stop... %s", loadingChars[i])
			i = (i + 1) % len(loadingChars)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
