package main

import (
	"log"
	"os"
	"sort"
	"sync"

	"github.com/DuGlaser/width-checker/model"
	"github.com/mxschmitt/playwright-go"
	"github.com/urfave/cli/v2"
)

var (
	max      uint
	min      uint
	interval uint
	url      string
	isOutput bool
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:        "max",
				Value:       1440,
				Usage:       "Max is the maximum width of the device to be checked.",
				Destination: &max,
			},
			&cli.UintFlag{
				Name:        "min",
				Value:       320,
				Usage:       "Min is the maximum width of the device to be checked.",
				Destination: &min,
			},
			&cli.UintFlag{
				Name:        "interval",
				Value:       50,
				Usage:       "Interval is the value of the change in width of the device.",
				Destination: &interval,
			},
			&cli.StringFlag{
				Name:        "url",
				Usage:       "Url is the url of the page to be checked.",
				Destination: &url,
				Required:    true,
			},
			&cli.BoolFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Destination: &isOutput,
			},
		},
		Action: func(c *cli.Context) error {
			deviceOption := &model.DeviceOption{
				Min:      int(min),
				Max:      int(max),
				Interval: int(interval),
			}

			pw, err := playwright.Run()
			if err != nil {
				log.Fatalf("could not start playwright: %v", err)
			}
			defer pw.Stop()

			var resultList model.ResultList
			var wg sync.WaitGroup

			browsers := []playwright.BrowserType{pw.Chromium, pw.Firefox}

			for _, bt := range browsers {
				wg.Add(1)
				go func(bt playwright.BrowserType) {
					defer wg.Done()
					WidthCheck(bt, &resultList, url, *deviceOption)
				}(bt)
			}

			wg.Wait()

			sort.Sort(model.ResultList(resultList))

			if isOutput {
				resultList.PrintResult()
				return nil
			}

			errorResultList := model.ResultList{}
			successResultList := model.ResultList{}

			for _, result := range resultList {
				if result.DeviceWidth == result.PageWidth {
					successResultList = append(successResultList, result)
				} else {
					errorResultList = append(errorResultList, result)
				}
			}

			if len(errorResultList) == 0 {
				errorResultList.PrintSuccess()
			} else {
				errorResultList.PrintError()
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func WidthCheck(bt playwright.BrowserType, resultList *model.ResultList, url string, opt model.DeviceOption) {
	browser, err := bt.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()
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
		defer context.Close()

		page, err := context.NewPage()
		if err != nil {
			log.Fatalf("could not create page: %v", err)
		}
		defer page.Close()

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

		result := model.Result{
			Browser:     bt.Name(),
			DeviceWidth: i,
			PageWidth:   pageWidth.(int),
		}
		*resultList = append(*resultList, result)
	}
}
