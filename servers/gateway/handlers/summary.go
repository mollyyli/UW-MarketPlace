package handlers

import (
	"encoding/json"
	"errors"
	"io"
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

const headerCORS = "Access-Control-Allow-Origin"
const corsAnyOrigin = "*"

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: add code and additional functions to do the following:
	// - Add an HTTP header to the response with the name
	//  `Access-Control-Allow-Origin` and a value of `*`. This will
	// 	allow cross-origin AJAX requests to your server.
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add(headerCORS, corsAnyOrigin)
	// - Get the `url` query string parameter value from the request.
	//   If not supplied, respond with an http.StatusBadRequest error.
	// - Call fetchHTML() to fetch the requested URL. See comments in that
	//   function for more details.
	url := r.URL.Query().Get("url")
	if len(url) == 0 {
		http.Error(w, `No url`, http.StatusBadRequest)
	}
	// - Call extractSummary() to extract the page summary meta-data,
	//   as directed in the assignment. See comments in that function
	//   for more details
	// - Close the response HTML stream so that you don't leak resources.
	// - Finally, respond with a JSON-encoded version of the PageSummary
	//   struct. That way the client can easily parse the JSON back into
	//   an object. Remember to tell the client that the response content
	//   type is JSON.
	summary, err := fetchHTML(url)
	if err != nil {
		// w.Write([]byte("error: " + err.Error()))
		if err.Error() == "404 not found" {
			http.Error(w, err.Error(), 404)
		}
		if err.Error() == "Content-Type is not HTML" {
			http.Error(w, err.Error(), 415)
		}
	} else {
		response, err := extractSummary(url, summary)

		json.NewEncoder(w).Encode(response)

		out, err := json.Marshal(response)
		if err == nil {
			r.FormValue(string(out))
		}

		// logger := log.New(os.Stdout, "htt: ", log.LstdFlags)
		// logger.Println(url)
	}
	// Helpful Links:
	// https://golang.org/pkg/net/http/#Request.FormValue
	// https://golang.org/pkg/net/http/#Error
	// https://golang.org/pkg/encoding/json/#NewEncoder

}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {

	resp, err := http.Get(pageURL)
	if err == nil {
		ctype := resp.Header.Get("Content-Type")
		if resp.StatusCode >= 400 {
			return nil, errors.New("404 not found")
		} else if !strings.HasPrefix(ctype, "text/html") {
			return nil, errors.New("Content-Type is not HTML")
		}
	} else {
		return nil, errors.New("404 not found")
	}
	return resp.Body, nil

}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.
	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary
	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/

	tokenizer := html.NewTokenizer(htmlStream)
	summary := new(PageSummary)

	var dict map[string]*PreviewImage
	dict = make(map[string]*PreviewImage)
	newIcon := new(PreviewImage)
	var arr []string

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
		}
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			if token.Data == "title" && len(summary.Title) == 0 {
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					summary.Title = tokenizer.Token().Data
				}
			}
			if token.Data == "meta" {
				var property string
				var content string
				var name string
				for _, class := range token.Attr {
					if class.Key == "property" {
						property = class.Val
					} else if class.Key == "content" {
						content = class.Val
					} else if class.Key == "name" {
						name = class.Val
					}
				}
				if strings.Contains(property, "og:image") {
					if property == "og:image" {
						image := new(PreviewImage)
						dict[content] = image
						dict[content].URL = content
						if !strings.Contains(dict[content].URL, "http") {
							dict[content].URL = handleAbsoluteURL(dict[content].URL, pageURL)
						}
						arr = append(arr, content)
					} else if property == "og:image:secure_url" {
						dict[arr[len(arr)-1]].SecureURL = content
					} else if property == "og:image:type" {
						dict[arr[len(arr)-1]].Type = content
					} else if property == "og:image:width" {
						widthInt, err := strconv.Atoi(content)
						if err == nil {
							dict[arr[len(arr)-1]].Width = widthInt
						}
					} else if property == "og:image:height" {
						heightInt, err := strconv.Atoi(content)
						if err == nil {
							dict[arr[len(arr)-1]].Height = heightInt
						}
					} else if property == "og:image:alt" {
						dict[arr[len(arr)-1]].Alt = content
					}
				} else if property == "og:type" {
					summary.Type = content
				} else if property == "og:url" {
					summary.URL = content
				} else if property == "og:title" {
					summary.Title = content
				} else if property == "og:site_name" {
					summary.SiteName = content
				} else if property == "og:description" {
					summary.Description = content
				} else if name == "author" {
					summary.Author = content
				} else if name == "keywords" {
					if strings.Contains(content, " ") {
						summary.Keywords = strings.Split(content, ", ")
					} else {
						summary.Keywords = strings.Split(content, ",")
					}
				} else if name == "description" && len(summary.Description) == 0 {
					summary.Description = content
				}
			}

			if "link" == token.Data {
				var rel string
				var icon string
				for _, class := range token.Attr {
					if class.Key == "rel" {
						rel = class.Key
						icon = class.Val
					}
				}
				for _, attribute := range token.Attr {
					if rel == "rel" && icon == "icon" {
						if attribute.Key == "href" {
							if !strings.Contains(attribute.Val, "http") {
								newIcon.URL = handleAbsoluteURL(attribute.Val, pageURL)
							} else {
								newIcon.URL = attribute.Val
							}
						}
						if attribute.Key == "sizes" && attribute.Val != "any" {
							sizes := strings.Split(attribute.Val, "x")
							if len(sizes) > 0 {
								height, err := strconv.Atoi(sizes[0])
								width, err := strconv.Atoi(sizes[1])
								if err != nil {
									return nil, err
								}
								newIcon.Height = height
								newIcon.Width = width
							}
						}
						if attribute.Key == "type" {
							newIcon.Type = attribute.Val
						}
						summary.Icon = newIcon
					}
				}
			}
		}
	}
	for _, link := range arr {
		summary.Images = append(summary.Images, dict[link])
	}
	return summary, nil
}

func handleAbsoluteURL(PageURL string, resourceURL string) string {
	URL, err := url.Parse(PageURL)
	if err == nil {
		base, err := url.Parse(resourceURL)
		if err == nil {
			resolveReferenceURL := base.ResolveReference(URL)
			return resolveReferenceURL.String()
		}
	}
	return ""
}
