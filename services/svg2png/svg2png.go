package svg2png

import (
	"context"
	"fmt"
	"gitmonitor/services/utils"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(html, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("data:text/html,"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(fmt.Sprintf("document.write(`%s`)", html)).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil

		}),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByQuery),
		// chromedp.FullScreenshot(res, 100),
		// chromedp.Screenshot()
	}
}

func GetImage(svg string) []byte {
	allocContext, cancel := chromedp.NewExecAllocator(context.Background(), append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.WindowSize(1920, 1080),
	)...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocContext)
	defer cancel()

	var buf []byte
	htmlTemplate := `<!DOCTYPE html>
	<html>
		<body>%s</body>
	</html>`
	html := fmt.Sprintf(htmlTemplate, svg)
	err := chromedp.Run(ctx, elementScreenshot(html, `svg`, &buf))
	utils.CheckErr("GetImageSVG2PNG", err)

	return buf
}
