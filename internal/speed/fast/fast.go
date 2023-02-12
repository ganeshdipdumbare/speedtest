package fast

import (
	"errors"
	"fmt"
	"log"

	"github.com/ganeshdipdumbare/gospeedtest/internal/speed"

	"github.com/playwright-community/playwright-go"
)

type fastSpeed struct {
	url string
}

func NewSpeedChecker() *fastSpeed {
	return &fastSpeed{
		url: "https://fast.com",
	}
}

func (f *fastSpeed) GetSpeed() (*speed.GetSpeedResp, error) {
	downloadSpeedResp := make(chan speed.NetSpeed)
	uploadSpeedResp := make(chan speed.NetSpeed)

	runOptions := &playwright.RunOptions{
		Verbose: false,
	}

	err := playwright.Install(runOptions)
	if err != nil {
		log.Fatalf("could not install playwright dependencies: %v", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	context, err := browser.NewContext()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	page, err := context.NewPage()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	// Navigate to fast.com
	_, err = page.Goto(f.url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	go func() {
		page.WaitForLoadState()

		waitOptions := playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(0),
		}

		waitAndEvalSpeedSelector(page, "div#speed-value.speed-results-container.succeeded", "#speed", downloadSpeedResp, waitOptions)

		err = page.Click("#show-more-details-link")
		if err != nil {
			uploadSpeedResp <- speed.NetSpeed{
				Err: err,
			}
		}

		waitAndEvalSpeedSelector(page, "div#extra-details-container.align-container.succeeded", "#upload", uploadSpeedResp, waitOptions)

		context.Close()
		browser.Close()
		err = pw.Stop()
		if err != nil {
			fmt.Printf("error while closing playwright run: %v\n", err)
		}
	}()

	return &speed.GetSpeedResp{
		DownloadSpeedChannel: downloadSpeedResp,
		UploadSpeedChannel:   uploadSpeedResp,
	}, nil
}

func waitAndEvalSpeedSelector(page playwright.Page, containerSelector, selectorPrefix string, resultCh chan speed.NetSpeed, waitOptions ...playwright.PageWaitForSelectorOptions) {
	defer close(resultCh)

	waitDone := make(chan bool)
	defer close(waitDone)

	go func() {
		_, err := page.WaitForSelector(containerSelector, waitOptions...)
		if err != nil {
			resultCh <- speed.NetSpeed{
				Err: err,
			}
		}
		waitDone <- true
	}()

	breakIt := false
	var (
		expression    = `element => element.textContent`
		valueSelector = selectorPrefix + `-value`
		unitsSelector = selectorPrefix + `-units`
	)

	for {
		select {
		case <-waitDone:
			resultCh <- getSpeedValueUnit(page, valueSelector, unitsSelector, expression)
			breakIt = true

		default:
			resultCh <- getSpeedValueUnit(page, valueSelector, unitsSelector, expression)
		}
		if breakIt {
			break
		}
	}
}

func getSpeedValueUnit(page playwright.Page, valueSelector, unitsSelector, expression string) speed.NetSpeed {
	response := speed.NetSpeed{}
	value, err := page.EvalOnSelector(valueSelector, expression)
	if err != nil {
		response.Err = err
		return response
	}

	unit, err := page.EvalOnSelector(unitsSelector, expression)
	if err != nil {
		response.Err = err
		return response
	}

	valueStr, ok := value.(string)
	if !ok {
		response.Err = errors.New("unable to convert speed to string")
		return response
	}

	unitStr, ok := unit.(string)
	if !ok {
		response.Err = errors.New("unable to convert unit to string")
		return response
	}

	response.Value = valueStr
	response.Unit = unitStr
	return response
}
