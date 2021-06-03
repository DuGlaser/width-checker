package model

import (
	"fmt"

	"github.com/fatih/color"
)

type Result struct {
	Browser     string
	DeviceWidth int
	PageWidth   int
}

func (r Result) infoLabel() string {
	return fmt.Sprintf("%-9v %-4vpx | ", r.Browser, r.DeviceWidth)
}

func (r Result) PrintSuccess() {
	browserOutput := r.infoLabel()
	fmt.Print(browserOutput)
	color.Green("SUCCESS")
}

func (r Result) PrintError() {
	browserOutput := r.infoLabel()
	resultOutput := fmt.Sprintf("ERROR: %vpx out of the device.", r.PageWidth-r.DeviceWidth)
	fmt.Print(browserOutput)
	color.Red(resultOutput)
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

func (r ResultList) PrintResult() {
	fmt.Println()

	for _, result := range r {
		if result.DeviceWidth == result.PageWidth {
			result.PrintSuccess()
		} else {
			result.PrintError()
		}
	}
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
