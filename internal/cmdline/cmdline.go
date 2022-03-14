package cmdline

import (
	"fmt"
	"time"
)

//顯示loading圈圈，status傳入要顯示在圈圈後的文字，delay傳入更新時間間隔，ch傳入結束訊號
func Spinner(status string, delay time.Duration, ch chan int) {
	for {
		select {
		case <-ch:
			fmt.Printf("\r                                              ")
			return
		default:
			for _, r := range `-\|/` {
				fmt.Printf("\r%c %s", r, status)
				time.Sleep(delay)
			}
		}
	}
}

//顯示進度(幾分之幾)，denominator傳入總數量(分母)，delay傳入更新時間間隔，in傳入目前進度數量(分子)，quit傳入結束訊號
func percentViewer(denominator int, delay time.Duration, in chan int, quit chan int) {
	var numerator int
	for {

		select {
		case <-quit:
			fmt.Printf("\r                                              ")
			return
		case numerator = <-in:
			fmt.Printf("\r>%d/%d", numerator, denominator)
		default:
			fmt.Printf("\r>%d/%d", numerator, denominator)
		}
	}
}
