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
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

func GetImage(svg string) []byte {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	var buf []byte
	htmlTemplate := `<!DOCTYPE html>
	<html>
		<body>%s</body>
	</html>`
	html := fmt.Sprintf(htmlTemplate, svg)
	err := chromedp.Run(ctx, elementScreenshot(html, `svg`, &buf))
	utils.CheckErr(err)

	return buf
}
