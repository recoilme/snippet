package snippet

import (
	"crypto/tls"
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
	Section     string `json:"section"`
	Tag         string `json:"tag"`
	Keywords    string `json:"keywords"`
	StatusCode  int    `json:"status_code"`
}

// SnippetFromReader return snippet from io.Reader
func SnippetFromReader(r io.ReadCloser) (*item, error) {
	defer r.Close()
	tokens := html.NewTokenizer(r)
	bodyFound := false
	it := &item{}
	tags := make([]string, 0, 1)
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
				}
			case "description":
				if tt == html.StartTagToken {
					descrText := tokens.Next()
					if descrText == html.TextToken {
						it.Description = strings.TrimSpace(tokens.Token().Data)
					}
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
							break
						}
						if strings.ToLower(attr.Val) == "keywords" {
							for _, attr := range tkn.Attr {
								if attr.Key == "content" {
									it.Keywords = strings.TrimSpace(attr.Val)
									break
								}
							}
							break
						}
					}
					if attr.Key == "property" {
						if strings.ToLower(attr.Val) == "article:section" {
							for _, attr := range tkn.Attr {
								if attr.Key == "content" {
									it.Section = strings.TrimSpace(attr.Val)
									break
								}
							}
						}
						if strings.ToLower(attr.Val) == "article:tag" {
							for _, attr := range tkn.Attr {
								if attr.Key == "content" {
									tags = append(tags, strings.TrimSpace(attr.Val))
									break
								}
							}
						}
					}
				}
			case "body":
				bodyFound = true
			}
		}

		if len(tags) > 0 {
			it.Tag = strings.Join(tags, ",")
		}
		if err || bodyFound {
			break
		}
	}
	return it, nil
}

// Snippet return snippet from url
// params: timeout (in seconds) and custom http headers
func Snippet(link string, timeout int, headers map[string]string) (*item, error) {
	r, statusCode, err := GetLinkReader(link, timeout, headers, 0)
	if err != nil {
		return &item{StatusCode: statusCode}, err
	}
	it, err := SnippetFromReader(r)
	it.StatusCode = statusCode
	return it, err
}

// getLinkReader return NopCloser from url
// params are timeout in seconds and headers
// dont foget to close reader
func GetLinkReader(link string, timeout int, headers map[string]string, maxSizeInt int) (io.ReadCloser, int, error) {
	// params
	maxSize := int64(maxSizeInt)
	if maxSizeInt <= 0 {
		maxSize = int64(300 * 1024) // 300 килобайт в байтах
	}
	defHeaders := make(map[string]string)
	defHeaders["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:65.0) Gecko/20100101 Firefox/65.0"
	defHeaders["Accept"] = "text/html,application/xhtml+xml,application/xml,application/rss+xml;q=0.9,image/webp,*/*;q=0.8"
	defHeaders["Accept-Language"] = "en-US;q=0.7,ru;q=0.3"
	if timeout == 0 {
		timeout = 10
	}
	// client
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: customTransport,
	}
	// request
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, 0, err
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
		return nil, 0, err
	}
	// return
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
		if err != nil {
			return nil, resp.StatusCode, err
		}
		limitedReader := io.LimitReader(utf8, maxSize)
		return io.NopCloser(limitedReader), resp.StatusCode, err
	}

	return nil, resp.StatusCode, fmt.Errorf("%s: %d", "Error, status code:", resp.StatusCode)
}
