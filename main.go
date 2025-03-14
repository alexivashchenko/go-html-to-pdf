package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {

	flag_chrome_user_data_dir := flag.String("chrome_user_data_dir", "", "Chrome user data directory.")
	flag_input_file := flag.String("input_file", "", "HTML file with badge code.")
	flag_output_file := flag.String("output_file", "", "File to output PDF.")
	flag_chrome := flag.String("chrome", "", "Path to Chrome executable.")

	flag_width := flag.Int("width", 400, "Badge Width,px.")
	flag_height := flag.Int("height", 600, "Badge Height,px.")
	flag.Parse()

	parentContext, parentContextCancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.ExecPath(*flag_chrome),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("lang", "en-US,en;q=0.9"),
		chromedp.UserDataDir(*flag_chrome_user_data_dir),
	)
	defer parentContextCancel()

	// ctx, cancel := chromedp.NewContext(parentContext, chromedp.WithDebugf(log.Printf))
	ctx, cancel := chromedp.NewContext(parentContext)
	defer cancel()

	html, err := os.ReadFile(*flag_input_file)
	if err != nil {
		fmt.Print(err)
	}

	pxInInch := float64(96)
	width := float64(*flag_width*2+30) / pxInInch
	height := float64(*flag_height+20) / pxInInch

	if err := chromedp.Run(ctx,

		chromedp.Navigate("about:blank"),

		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, cancel := context.WithCancel(ctx)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(1)
			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					cancel()
					wg.Done()
				}
			})

			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			if err := page.SetDocumentContent(frameTree.Frame.ID, string(html)).Do(ctx); err != nil {
				return err
			}
			wg.Wait()
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(width).
				WithPaperHeight(height).
				Do(ctx)

			if err != nil {
				return err
			}
			return ioutil.WriteFile(*flag_output_file, buf, 0644)
		}),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}
