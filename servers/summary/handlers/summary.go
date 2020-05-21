package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	pageURL := r.FormValue("url")
	if pageURL == "" {
		http.Error(w, "No URL was provided.", http.StatusBadRequest)
		return
	}
	body, err := fetchHTML(pageURL)
	log.Println(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	summary, err := extractSummary(pageURL, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body.Close()
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(summary); err != nil {
		fmt.Printf("error encoding struct into JSON: %v\n", err)
		return
	}

}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(string(resp.StatusCode))
	}
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, errors.New(ctype)
	}
	return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	tokenizer := html.NewTokenizer(htmlStream)
	summary := new(PageSummary)
	icon := new(PreviewImage)
	images := []*PreviewImage{}
	index := strings.LastIndex(pageURL, "/")
	comIndex := strings.LastIndex(pageURL, ".com")
	if index > comIndex {
		pageURL = pageURL[:index]
	}

	descFound := false
	titleFound := false

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				return summary, nil
			}
			return nil, err
		}

		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			if "body" == token.Data {
				break
			}
			if "meta" == token.Data {

				key, value := extractMetaData(token)

				switch key {
				case "og:type":
					summary.Type = value
				case "og:url":
					summary.URL = value
				case "og:title":
					if !titleFound {
						summary.Title = value
						titleFound = true
					}
					summary.Title = value
				case "og:site_name":
					summary.SiteName = value
				case "og:description":
					summary.Description = value
					descFound = true
				case "description":
					if !descFound {
						summary.Description = value
					}
				case "author":
					summary.Author = value
				case "keywords":
					summary.Keywords = getKeywords(value)
				case "og:image":
					images = append(images, createNewPreviewImage(pageURL, value))
				case "og:image:secure_url":
					images[len(images)-1].SecureURL = value
				case "og:image:type":
					images[len(images)-1].Type = value
				case "og:image:width":
					images[len(images)-1].Width, _ = strconv.Atoi(value)
				case "og:image:height":
					images[len(images)-1].Height, _ = strconv.Atoi(value)
				case "og:image:alt":
					images[len(images)-1].Alt = value
				}

			}
			if "title" == token.Data {
				if !titleFound {
					tokenType = tokenizer.Next()
					if tokenType == html.TextToken {
						summary.Title = tokenizer.Token().Data
					}
				}
			}
			if "link" == token.Data {
				href, mediaType, sizes, isIcon := "", "", "", false
				for _, attr := range token.Attr {
					if attr.Key == "rel" && attr.Val == "icon" {
						isIcon = true
					}
					if attr.Key == "href" {
						href = attr.Val
					}
					if attr.Key == "sizes" {
						sizes = attr.Val
					}
					if attr.Key == "type" {
						mediaType = attr.Val
					}
				}
				if isIcon {
					testURL, _ := url.Parse(href)
					if !testURL.IsAbs() {
						icon.URL = pageURL + href
					} else {
						icon.URL = href
					}
					if mediaType != "" {
						icon.Type = mediaType
					}
					if sizes != "" && sizes != "any" {
						dimensions := strings.Split(sizes, "x")
						icon.Height, _ = strconv.Atoi(dimensions[0])
						icon.Width, _ = strconv.Atoi(dimensions[1])
					}
				}
			}
		}

		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if "head" == token.Data {
				break
			}
		}
	}
	if icon.URL != "" {
		summary.Icon = icon
	}
	if len(images) > 0 {
		summary.Images = images
	}
	return summary, nil
}

func extractMetaData(token html.Token) (string, string) {
	key, value := "", ""
	for _, attr := range token.Attr {
		if attr.Key == "property" || attr.Key == "name" {
			key = attr.Val
		}
		if attr.Key == "content" {
			value = attr.Val
		}
	}
	return key, value
}

func getKeywords(currList string) []string {
	modifiedList := strings.Split(currList, ",")
	for i, keyword := range modifiedList {
		modifiedList[i] = strings.TrimSpace(keyword)
	}
	return modifiedList
}

func createNewPreviewImage(pageURL string, imageURL string) *PreviewImage {
	image := new(PreviewImage)
	testURL, _ := url.Parse(imageURL)
	if !testURL.IsAbs() {
		image.URL = pageURL + imageURL
	} else {
		image.URL = imageURL
	}
	return image
}
