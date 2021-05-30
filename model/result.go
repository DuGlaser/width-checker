package model

import (
	"fmt"
)

type Result struct {
	Browser     string
	DeviceWidth int
	PageWidth   int
}

func (r *Result) PrintError() {
	browserOutput := fmt.Sprintf("[ %-9v]", r.Browser)
	resultOutput := fmt.Sprintf(" deviceWidth: %-4v  pageWidth: %-4v", r.DeviceWidth, r.PageWidth)
	fmt.Print(browserOutput)
	fmt.Println(resultOutput)
}

type ResultList []Result

func (r ResultList) Len() int      { return len(r) }
func (r ResultList) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r ResultList) Less(i, j int) bool {
	if r[i].DeviceWidth > r[j].DeviceWidth {
		return true
	}
	if r[i].DeviceWidth < r[j].DeviceWidth {
		return false
	}
	return r[i].Browser < r[j].Browser
}

func (r ResultList) PrintError() {
	fmt.Println()
	fmt.Println()
	fmt.Println("===== ERROR =====")
	fmt.Println()

	for _, result := range r {
		result.PrintError()
	}
}

func (r ResultList) PrintSuccess() {
	fmt.Println()
	fmt.Println("SUCCESS!")
}
