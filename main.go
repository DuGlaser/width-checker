package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/DuGlaser/width-checker/model"
	"github.com/cheggaaa/pb/v3"
	"github.com/mxschmitt/playwright-go"
)

func WidthCheck(bt playwright.BrowserType, resultList chan model.Result, url string, opt model.DeviceOption) {
	browser, err := bt.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	for i := opt.Min; i <= opt.Max; i += opt.Interval {
		context, err := browser.NewContext(playwright.BrowserNewContextOptions{
			IsMobile: playwright.Bool(false),
			Viewport: &playwright.BrowserNewContextOptionsViewport{
				Width:  playwright.Int(i),
				Height: playwright.Int(i),
			},
		})
		if err != nil {
			log.Fatalf("could not create context: %v", err)
		}

		page, err := context.NewPage()
		if err != nil {
			log.Fatalf("could not create page: %v", err)
		}

		_, err = page.Goto(
			url,
			playwright.PageGotoOptions{
				WaitUntil: playwright.WaitUntilStateNetworkidle,
			},
		)
		if err != nil {
			log.Fatalf("could not goto: %v", err)
		}

		handle, err := page.EvaluateHandle("document.querySelector('body')")
		if err != nil {
			log.Fatalf("could not acquire JSHandle: %v\n", err)
		}

		r, err := handle.(playwright.JSHandle).GetProperty("scrollWidth")
		if err != nil {
			log.Fatalf("could not get scrollWidth: %v\n", err)
		}
		pageWidth, _ := r.JSONValue()

		result := &model.Result{
			Browser:     bt.Name(),
			DeviceWidth: i,
			PageWidth:   pageWidth.(int),
		}
		resultList <- *result

		context.Close()
	}

	browser.Close()
}

func main() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	url := "http://example.com"
	deviceOption := &model.DeviceOption{
		Min:      320,
		Max:      1440,
		Interval: 200,
	}

	count := deviceOption.Pattern() * 2

	ch := make(chan model.Result, count)

	go WidthCheck(pw.Chromium, ch, url, *deviceOption)
	go WidthCheck(pw.Firefox, ch, url, *deviceOption)

	fmt.Println()

	errorResultList := model.ResultList{}

	bar := pb.StartNew(count)
	i := 0

	for result := range ch {
		i++
		bar.Increment()

		if result.DeviceWidth != result.PageWidth {
			errorResultList = append(errorResultList, result)
		}

		if i >= count {
			break
		}
	}

	bar.Finish()

	if len(errorResultList) == 0 {
		errorResultList.PrintSuccess()
	} else {
		sort.Sort(model.ResultList(errorResultList))
		errorResultList.PrintError()
	}

	pw.Stop()
}
