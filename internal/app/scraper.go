package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Aegon95/home24-webscraper/pkg/models"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
)

type AnalyzeWebService interface {
	Scraper(siteUrl string) (*models.WebStats, error)
}

type analyzer struct {
	client *http.Client
	logger *zap.SugaredLogger
}

func NewAnalyzeWebService(client *http.Client, logger *zap.SugaredLogger) AnalyzeWebService {
	return &analyzer{
		logger: logger,
		client: client,
	}
}

// Scraper function scrapes the website and calculates the stats
func (a *analyzer) Scraper(siteUrl string) (*models.WebStats, error) {

	// make Get request to the entered URL
	res, err := a.client.Get(siteUrl)

	if err != nil || res.StatusCode != http.StatusOK {
		a.logger.Errorf("Error occurred while making Get Request to %s", siteUrl)
		return nil, errors.New("Cannot reach the entered URL")
	}

	// parse the response body
	node, err := html.Parse(res.Body)

	// create result webstats data
	webCount := &models.WebStats{}


	// create context to pass data
	ctx := context.WithValue(context.Background(), "url", siteUrl)

	// set default values
	webCount.HTMLVersion = "5.0"
	webCount.HasLogin = false

	if u, ok := ctx.Value("url").(string); ok {
		parsedUrl, _ := url.Parse(u)
		ctx = context.WithValue(ctx, "mainUrl", parsedUrl)
	}

	requestWg := &sync.WaitGroup{}

	ctx = context.WithValue(ctx, "requestWg", requestWg)
	ctx = context.WithValue(ctx, "webCount", webCount)

	var f func(node *html.Node)

	// loop over the node and its child nodes
	f = func(node *html.Node) {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			a.processNode(ctx, c)
			f(c)
		}
	}

	f(node)

	defer res.Body.Close()

	// wait for all go routines to complete
	requestWg.Wait()

	return webCount, nil
}

// processNode processes a node based on its type
func (a *analyzer) processNode(ctx context.Context, node *html.Node) error {
	if node == nil {
		return nil
	}

	if node.Type == html.DoctypeNode{
		if u, ok := ctx.Value("webCount").(*models.WebStats); ok {
			u.HTMLVersion = getHTMLVersion(node.Attr)
		}

	}else{
		if err := a.processElementNode(ctx,node); err != nil {
			return err
		}
	}
	return nil
}

// processElementNode updates the stats based on type
func (a *analyzer) processElementNode(ctx context.Context, node *html.Node) error {

	if stats, ok := ctx.Value("webCount").(*models.WebStats); ok {
		switch node.Data {
		case "title":
			if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
				stats.Title = node.FirstChild.Data
			}

		case "h1":
			stats.H1Count++
		case "h2":
			stats.H2Count++
		case "h3":
			stats.H3Count++
		case "h4":
			stats.H4Count++
		case "h5":
			stats.H5Count++
		case "h6":
			stats.H6Count++
		case "a":
			if attr := findAttributeByKey(node.Attr, "href"); attr != nil {
				if err := a.processLink(ctx,attr.Val); err != nil {
					fmt.Printf("error analyzing anchor link (%v): %v", attr.Val, err)
				}

			}

		case "form":
			if attr := findAttributeByKey(node.Attr, "action"); attr != nil {
				if err := a.processLink(ctx, attr.Val); err != nil {
					fmt.Printf("error analyzing form link (%v): %v", attr.Val, err)
				}
			}

		case "input":
			if checkAttributeEquals(node, "type", "password") {
				stats.HasLogin = true
			}

		}
	}


	return nil
}

// processLink checks if link is external or internal
func (a *analyzer) processLink(ctx context.Context, link string) error {
	stats := &models.WebStats{}
	if stats, ok := ctx.Value("webCount").(*models.WebStats); ok {
		if len(link) == 0 {
			stats.InternalLinksCount++

		}
	}

	var pageHost string
	if parsedUrl, ok := ctx.Value("mainUrl").(*url.URL); ok {
		pageHost = parsedUrl.Host
	}

	parsed, err := url.Parse(link)
	if err != nil {
		a.logger.Errorf("Error occurred while parsing the link %s", link)
		return err
	}

	if !parsed.IsAbs() || parsed.Host == pageHost {
		stats.InternalLinksCount++
		err := a.checkLinkAccessible(ctx,link,false)
		if err != nil{
			return err
		}

	} else {
		stats.ExternalLinksCount++
		err := a.checkLinkAccessible(ctx,link,true)
		if err != nil{
			return err
		}
	}

	return nil
}

// checkLinkAccessible checks if link is not accessible
func (a *analyzer) checkLinkAccessible(ctx context.Context, link string, isAbs bool) error {

	stats, ok := ctx.Value("webCount").(*models.WebStats)
	requestWg, ok := ctx.Value("requestWg").(*sync.WaitGroup)
	mainUrl, ok := ctx.Value("mainUrl").(*url.URL)
	if !ok{
		return errors.New("Error while reading webCount")
	}
	requestWg.Add(1)
	go func(link string, homeUrl *url.URL) {
		defer requestWg.Done()

		if !isAbs {
			var currentPage string
			if homeUrl != nil {
				currentPage = homeUrl.String()
			}
			link = currentPage + "/" + link
		}

		resp, err := http.Head(link)
		if err != nil || resp == nil || resp.StatusCode >= 400 || resp.StatusCode < 200 {
			atomic.AddInt32(&stats.InaccessibleLinksCount, 1)
		}
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}(link, mainUrl)

	return nil
}

// helper functions

func getHTMLVersion(doctypeAttrs []html.Attribute) string {
	if attr := findAttributeByKey(doctypeAttrs, "public"); attr != nil && strings.Contains(attr.Val, "4.0") {
		return "4.0"
	} else {
		return "5.0"
	}
}

func findAttributeByKey(attrs []html.Attribute, key string) *html.Attribute {
	for i := range attrs {
		if attrs[i].Key == key {
			return &attrs[i]
		}
	}
	return nil
}

func checkAttributeEquals(node *html.Node, key, value string) bool {
	if attr := findAttributeByKey(node.Attr, key); attr != nil && attr.Val == value {
		return true
	}
	return false
}
