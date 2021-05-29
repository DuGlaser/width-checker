package main

import (
	"fmt"
	"log"

	"github.com/morikuni/aec"
	"github.com/mxschmitt/playwright-go"
)

type DeviceOption struct {
	min      int
	max      int
	interval int
}

type Result struct {
	browser     string
	width       int
	scrollWidth interface{}
}

var (
	green = aec.Color8BitF(aec.NewRGB8Bit(64, 255, 64))
	red   = aec.EmptyBuilder.LightRedF().ANSI
)

func widthCheck(bt playwright.BrowserType, url string, opt DeviceOption) {
	browser, err := bt.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	for i := opt.min; i <= opt.max; i += opt.interval {
		context, err := browser.NewContext(playwright.BrowserNewContextOptions{
			IsMobile: playwright.Bool(false),
			Viewport: &playwright.BrowserNewContextOptionsViewport{
				Width:  playwright.Int(i),
				Height: playwright.Int(i),
			},
		})
		defer context.Close()
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
		scrollWidth, _ := r.JSONValue()

		result := &Result{
			browser:     bt.Name(),
			width:       i,
			scrollWidth: scrollWidth,
		}

		browserOutput := fmt.Sprintf("[ %-9v]", result.browser)
		resultOutput := fmt.Sprintf(" width: %-4v  scrollWidth: %-4v", result.width, result.scrollWidth)

		fmt.Print(green.Apply(browserOutput))
		if result.width == result.scrollWidth {
			fmt.Println(green.Apply(resultOutput))
		} else {
			fmt.Println(red.Apply(resultOutput))
		}
	}
}

func main() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	url := ""

	ch1 := make(chan bool)
	ch2 := make(chan bool)

	go func(pw *playwright.Playwright) {
		widthCheck(pw.Chromium, url, *&DeviceOption{
			min:      320,
			max:      1440,
			interval: 50,
		})

		ch1 <- true
	}(pw)

	go func(pw *playwright.Playwright) {
		widthCheck(pw.Firefox, url, *&DeviceOption{
			min:      320,
			max:      1440,
			interval: 50,
		})

		ch2 <- true
	}(pw)

	<-ch1
	<-ch2
}
