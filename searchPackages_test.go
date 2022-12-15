package bdd_godog

import (
	"context"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/playwright-community/playwright-go"
	"log"
	"strings"
)

const godogCtxBrowserKey = "godogCtxBrowserKey"
const godogCtxPageKey = "godogCtxPageKey"

func iEnterPackageNameInTheSearch(ctx context.Context, packageName string) error {
	page, ok := ctx.Value(godogCtxPageKey).(playwright.Page)
	if !ok {
		return errors.New("page instance not found")
	}

	locator, err := page.Locator("xpath=//input[@id='AutoComplete']")
	if err != nil {
		return err
	}

	err = locator.Fill(packageName)
	if err != nil {
		return err
	}

	return nil
}

func iNavigateToHttpspkggodevByClickingPackagesOnMenu(ctx context.Context) error {
	page, ok := ctx.Value(godogCtxPageKey).(playwright.Page)
	if !ok {
		return errors.New("page instance not found")
	}

	locator, err := page.Locator("xpath=//header/div[1]/nav[1]/div[1]/ul[1]/li[4]/a[1]")
	if err != nil {
		return err
	}

	err = locator.Click()
	if err != nil {
		return err
	}

	pageUrl := page.URL()
	if pageUrl != "https://pkg.go.dev/" {
		return errors.New("failed to navigate to packages url")
	}

	return nil
}

func iPressSearchButton(ctx context.Context) error {
	page, ok := ctx.Value(godogCtxPageKey).(playwright.Page)
	if !ok {
		return errors.New("page instance not found")
	}

	locator, err := page.Locator("xpath=//button[contains(text(),'Search')]")
	if err != nil {
		return err
	}

	err = locator.Click()
	if err != nil {
		return err
	}
	return nil
}

func iShouldSeeASearchPageWithPackage(ctx context.Context, packageUrl string) error {
	page, ok := ctx.Value(godogCtxPageKey).(playwright.Page)
	if !ok {
		return errors.New("page instance not found")
	}

	var timeOut float64 = 2000
	_, err := page.WaitForSelector(".SearchSnippet-headerContainer", playwright.PageWaitForSelectorOptions{
		Timeout: &timeOut,
	})

	query, err := page.QuerySelectorAll(".SearchSnippet-headerContainer")
	if err != nil {
		return err
	}

	isAvailable := false
	for _, element := range query {
		titleElement, err := element.QuerySelector("h2> a")
		if err != nil {
			return err
		}

		title, err := titleElement.TextContent()
		if err != nil {
			return err
		}

		if strings.Contains(title, packageUrl) {
			isAvailable = true
			break
		}
	}

	if !isAvailable {
		return errors.New(fmt.Sprintf("package %s is not available", packageUrl))
	}

	return nil
}

func iVisit(ctx context.Context, url string) (context.Context, error) {
	browser, ok := ctx.Value(godogCtxBrowserKey).(playwright.Browser)
	if !ok {
		return ctx, errors.New("failed to get browser instance")
	}

	browserContext, err := browser.NewContext()
	if err != nil {
		return ctx, err
	}

	page, err := browserContext.NewPage()
	if err != nil {
		return ctx, err
	}

	_, err = page.Goto(url)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, godogCtxPageKey, page), nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("failed to instanstiate Playwright instance: %v ", err)
	}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		option := playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(false),
			Channel:  playwright.String("chrome"),
		}

		browser, err := pw.Chromium.Launch(option)
		if err != nil {
			return ctx, err
		}

		return context.WithValue(ctx, godogCtxBrowserKey, browser), nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		browser, ok := ctx.Value(godogCtxBrowserKey).(playwright.Browser)
		if !ok {
			return ctx, errors.New("failed to get browser instance")
		}
		err = browser.Close()
		if err != nil {
			return ctx, errors.New("failed to close browser instance")
		}

		return ctx, nil
	})

	ctx.Step(`^I enter "([^"]*)" package name in the search$`, iEnterPackageNameInTheSearch)
	ctx.Step(`^I navigate to https:\/\/pkg\.go\.dev by clicking packages on menu$`, iNavigateToHttpspkggodevByClickingPackagesOnMenu)
	ctx.Step(`^I press search button$`, iPressSearchButton)
	ctx.Step(`^I should see a search page with "([^"]*)" package$`, iShouldSeeASearchPageWithPackage)
	ctx.Step(`^I visit "([^"]*)"$`, iVisit)
}
