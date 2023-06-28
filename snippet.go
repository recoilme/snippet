package snippet

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

type item struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func snippet(link string, timeout int, headers map[string]string) (*item, error) {
	r, err := getLinkReader(link, timeout, headers)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	tokens := html.NewTokenizer(r)
	titleFind := false
	descriptionFind := false
	it := &item{}
	for {
		tt := tokens.Next()
		err := false
		switch tt {
		case html.ErrorToken:
			err = true
		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			tkn := tokens.Token()
			switch tkn.Data {
			case "title":
				if tt == html.StartTagToken {
					titleText := tokens.Next()
					if titleText == html.TextToken {
						it.Title = strings.TrimSpace(tokens.Token().Data)
					}
					titleFind = true
				}
			case "description":
				if tt == html.StartTagToken {
					descrText := tokens.Next()
					if descrText == html.TextToken {
						it.Description = strings.TrimSpace(tokens.Token().Data)
					}
					//descriptionFind = true//? search untip meta or stop?
				}
			case "meta":
				for _, attr := range tkn.Attr {
					if attr.Key == "name" {
						if strings.ToLower(attr.Val) == "description" {
							for _, attr := range tkn.Attr {
								if attr.Key == "content" {
									it.Description = strings.TrimSpace(attr.Val)
									break
								}
							}
							descriptionFind = true
							break
						}
					}
				}
			}
		}

		if (titleFind && descriptionFind) || err {
			break
		}
	}
	r.Close()
	return it, nil
}

// getLinkReader return NopCloser from url
// params are timeout in seconds and headers
// dont foget to close reader
func getLinkReader(link string, timeout int, headers map[string]string) (io.ReadCloser, error) {
	// params
	defHeaders := make(map[string]string)
	defHeaders["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:65.0) Gecko/20100101 Firefox/65.0"
	defHeaders["Accept"] = "text/html,application/xhtml+xml,application/xml,application/rss+xml;q=0.9,image/webp,*/*;q=0.8"
	defHeaders["Accept-Language"] = "en-US;q=0.7,en;q=0.3"
	if timeout == 0 {
		timeout = 10
	}
	// client
	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	// request
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}
	// headers
	for k, v := range defHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	// response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// return
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
		if err != nil {
			return nil, err
		}
		return io.NopCloser(utf8), err
	}

	return nil, fmt.Errorf("%s: %d", "Error, status code:", resp.StatusCode)
}
