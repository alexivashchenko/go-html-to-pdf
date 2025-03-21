package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {

	log.Println("START")
	start := time.Now()

	// gogo --input_files='i-1.html,i-2.html,i-3.html' --output_files='o-1.pdf,o-2.pdf,o-3.pdf'
	// powershell Start-Process -FilePath "go-html-to-pdf.exe" -NoNewWindow -Wait -ArgumentList '-input_file="W:\SERVER\files\accreditation-app\storage\f8467ca7-95b5-4555-ad6e-be39da06677c.html" -output_file="W:\SERVER\files\accreditation-app\storage\f8467ca7-95b5-4555-ad6e-be39da06677c.pdf" -width=400 -height=600 -chrome="C:\Users\lex\.cache\puppeteer\chrome\win64-126.0.6478.182\chrome-win64\chrome.exe" -chrome_user_data_dir="/"
	/*
		-input_file="W:\SERVER\files\accreditation-app\storage\f8467ca7-95b5-4555-ad6e-be39da06677c.html"
		-output_file="W:\SERVER\files\acgo build main.gocreditation-app\storage\f8467ca7-95b5-4555-ad6e-be39da06677c.pdf"
		gogo

		-input_files="W:\SERVER\files\accreditation-app\storage\0c7f9705-02e4-4625-b834-6c2dfec9e8ef.html,W:\SERVER\files\accreditation-app\storage\9bf4f766-6934-4909-af66-0a5279408355.html,W:\SERVER\files\accreditation-app\storage\805f2ebf-fc49-4dfb-8c72-d1cd626e21d8.html,W:\SERVER\files\accreditation-app\storage\e57d4ef3-4001-4755-9716-2248621bc8c5.html,W:\SERVER\files\accreditation-app\storage\f8467ca7-95b5-4555-ad6e-be39da06677c.html"
		-output_files="W:\SERVER\files\accreditation-app\storage\0c7f9705-02e4-4625-b834-6c2dfec9e8ef.pdf,W:\SERVER\files\accreditation-app\storage\9bf4f766-6934-4909-af66-0a5279408355.pdf,W:\SERVER\files\accreditation-app\storage\805f2ebf-fc49-4dfb-8c72-d1cd626e21d8.pdf,W:\SERVER\files\accreditation-app\storage\e57d4ef3-4001-4755-9716-2248621bc8c5.pdf,W:\SERVER\files\accreditation-app\storage\f8467ca7-95b5-4555-ad6e-be39da06677c.pdf"

		-widths="400,400,400,400,400"
		-heights="600,600,600,600,600"
		-chrome="C:\Users\lex\.cache\puppeteer\chrome\win64-126.0.6478.182\chrome-win64\chrome.exe"
		-chrome_user_data_dir="/"


		gogo -input_files="W:\SERVER\files\accreditation-app\storage\0c7f9705-02e4-4625-b834-6c2dfec9e8ef.html,W:\SERVER\files\accreditation-app\storage\9bf4f766-6934-4909-af66-0a5279408355.html,W:\SERVER\files\accreditation-app\storage\805f2ebf-fc49-4dfb-8c72-d1cd626e21d8.html,W:\SERVER\files\accreditation-app\storage\e57d4ef3-4001-4755-9716-2248621bc8c5.html,W:\SERVER\files\accreditation-app\storage\f8467ca7-95b5-4555-ad6e-be39da06677c.html" -output_files="W:\SERVER\files\accreditation-app\storage\0c7f9705-02e4-4625-b834-6c2dfec9e8ef.pdf,W:\SERVER\files\accreditation-app\storage\9bf4f766-6934-4909-af66-0a5279408355.pdf,W:\SERVER\files\accreditation-app\storage\805f2ebf-fc49-4dfb-8c72-d1cd626e21d8.pdf,W:\SERVER\files\accreditation-app\storage\e57d4ef3-4001-4755-9716-2248621bc8c5.pdf,W:\SERVER\files\accreditation-app\storage\f8467ca7-95b5-4555-ad6e-be39da06677c.pdf" -widths="400,400,400,400,400" -heights="600,600,600,600,600" -chrome="C:\Users\lex\.cache\puppeteer\chrome\win64-126.0.6478.182\chrome-win64\chrome.exe" -chrome_user_data_dir="/"
		gogo -input_files="W:\SERVER\files\accreditation-app\storage\1da83f4f-5b65-4695-a399-32591cb5162a.html" -output_files="W:\SERVER\files\accreditation-app\storage\1da83f4f-5b65-4695-a399-32591cb5162a.pdf" -widths="400" -heights="600" -chrome="C:\Users\lex\.cache\puppeteer\chrome\win64-126.0.6478.182\chrome-win64\chrome.exe" -chrome_user_data_dir="/"
	*/

	flag_chrome_user_data_dir := flag.String("chrome_user_data_dir", "", "Chrome user data directory.")
	flag_input_files := flag.String("input_files", "", "HTML file with badge code.") // --input_files='i-1.html,i-2.html,i-3.html'
	flag_output_files := flag.String("output_files", "", "File to output PDF.")      // --output_files='o-1.pdf,o-2.pdf,o-3.pdf'
	flag_chrome := flag.String("chrome", "", "Path to Chrome executable.")
	flag_widths := flag.String("widths", "", "Badge Width,px.")
	flag_heights := flag.String("heights", "", "Badge Height,px.")
	flag.Parse()

	inputFiles := strings.Split(*flag_input_files, ",")
	outputFiles := strings.Split(*flag_output_files, ",")
	widths := strings.Split(*flag_widths, ",")
	heights := strings.Split(*flag_heights, ",")

	wg := &sync.WaitGroup{}

	for fileIndex, inputFile := range inputFiles {

		width, _ := strconv.Atoi(widths[fileIndex])
		height, _ := strconv.Atoi(heights[fileIndex])

		wg.Add(1)
		go convertFile(
			wg,
			inputFile,
			outputFiles[fileIndex],
			width,
			height,
			*flag_chrome_user_data_dir,
			*flag_chrome,
		)
	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("TIME %s", elapsed)

	fmt.Println("Done")
}

func convertFile(
	wg *sync.WaitGroup,
	input_file string,
	output_file string,
	flag_width int,
	flag_height int,
	flag_chrome_user_data_dir string,
	flag_chrome string,
) {

	defer wg.Done()

	fmt.Println("input_file:", input_file)
	fmt.Println("output_file:", output_file)

	parentContext, parentContextCancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.ExecPath(flag_chrome),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("lang", "en-US,en;q=0.9"),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("lang", "en-US,en;q=0.9"),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("single-process", true), // Faster for single conversions
		chromedp.UserDataDir(flag_chrome_user_data_dir),
	)
	defer parentContextCancel()

	// ctx, cancel := chromedp.NewContext(parentContext, chromedp.WithDebugf(log.Printf))
	ctx, cancel := chromedp.NewContext(parentContext)
	defer cancel()

	html, err := os.ReadFile(input_file)
	if err != nil {
		fmt.Print(err)
	}

	pxInInch := float64(96)
	width := float64(flag_width*2+30) / pxInInch
	height := float64(flag_height+20) / pxInInch

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
			return ioutil.WriteFile(output_file, buf, 0644)
		}),
	); err != nil {
		log.Fatal(err)
	}
}
