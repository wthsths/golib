package routing

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	gl_http "github.com/payports/golib/http"
)

type proxyClient struct {
	routeTable *RouteTable
	routeUrl   string
	httpCli    *http.Client
	onErr      func(error)
	onReqRead  func([]byte)
	onResRead  func([]byte)
}

// NewProxyClient creates a new proxy client instance.
// Underlying HandleRequestAndRedirect method can be registered as a handler function.
// Handler function will redirect incoming request to the routeUrl.
//
// Since ListenAndServe is blocking function, internally occuring event data can be captured via following hooks:
//
// onErr: Can be registered to receive internal errors.
//
// onReqRead: Can be registered to get incoming request body.
//
// onResRead: Can be registered to get outgoing response body.
func NewProxyClient(routeTable *RouteTable, routeUrl string, httpCli *http.Client, onErr func(error), onReqRead func([]byte), onResRead func([]byte)) *proxyClient {
	return &proxyClient{
		routeTable: routeTable,
		routeUrl:   routeUrl,
		httpCli:    httpCli,
		onErr:      onErr,
		onReqRead:  onReqRead,
		onResRead:  onResRead,
	}
}

// HandleRequestAndRedirect can be registered to http.Handle() for redirecting requests to desired url.
func (pc *proxyClient) HandleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.RequestURI()

	var err error
	isAllowed := false

	for _, e := range pc.routeTable.routeRules {
		if e.method != r.Method {
			continue
		}
		if e.regexp.MatchString(uri) {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("path is not allowed: %s", uri))
		}
		writtenRes, err := gl_http.WriteCustomJsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"message": "unauthorized call",
		})
		if err != nil && pc.onErr != nil {
			pc.onErr(fmt.Errorf("write response error: %w", err))
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes)
		}
		return
	}

	redirectUrl := pc.routeUrl + r.URL.RequestURI()
	parsedRedirectUrl, err := url.Parse(redirectUrl)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("unable to parse URL: '%s' error: %w", redirectUrl, err))
		}
		writtenRes, err := gl_http.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil && pc.onErr != nil {
			pc.onErr(fmt.Errorf("write response error: %w", err))
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes)
		}
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("error reading request body: %w", err))
		}
		writtenRes, err := gl_http.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil && pc.onErr != nil {
			pc.onErr(fmt.Errorf("write response error: %w", err))
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes)
		}
		return
	}
	defer r.Body.Close()

	if pc.onReqRead != nil {
		pc.onReqRead(reqBytes)
	}

	buffer := bytes.NewBuffer(reqBytes)
	nopCloser := ioutil.NopCloser(buffer)

	httpReq := &http.Request{
		Method: r.Method,
		URL:    parsedRedirectUrl,
		Header: r.Header,
		Body:   nopCloser,
	}

	res, err := pc.httpCli.Do(httpReq)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("error executing http request: %w", err))
		}
		writtenRes, err := gl_http.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil && pc.onErr != nil {
			pc.onErr(fmt.Errorf("write response error: %w", err))
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes)
		}
		return
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("error reading response payload: %w", err))
		}
		writtenRes, err := gl_http.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil && pc.onErr != nil {
			pc.onErr(fmt.Errorf("write response error: %w", err))
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes)
		}
		return
	}
	defer res.Body.Close()

	for k, v := range res.Header {
		for i := 0; i < len(v); i++ {
			w.Header().Add(k, v[i])
		}
	}

	_, err = w.Write(resBytes)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("error writing server response for client: %w", err))
		}
		writtenRes, err := gl_http.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil && pc.onErr != nil {
			pc.onErr(fmt.Errorf("write response error: %w", err))
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes)
		}
		return
	}

	/* Set status code last, or other header values and body will be lost. */
	if res.StatusCode != http.StatusOK {
		w.WriteHeader(res.StatusCode)
	}

	if res.Header.Get("Content-Encoding") == "gzip" {
		reader := bytes.NewReader(resBytes)
		gzipReader, err := gzip.NewReader(reader)
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(fmt.Errorf("error creating gzip reader: %w", err))
			}
			writtenRes, err := gl_http.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "internal error",
			})
			if err != nil && pc.onErr != nil {
				pc.onErr(fmt.Errorf("write response error: %w", err))
			}
			if pc.onResRead != nil {
				pc.onResRead(writtenRes)
			}
			return
		}
		/* Modifying resBytes for logging decompressed content AFTER we've written the response body. */
		resBytes, err = ioutil.ReadAll(gzipReader)
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(fmt.Errorf("error reading from gzip reader: %w", err))
			}
			writtenRes, err := gl_http.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "internal error",
			})
			if err != nil && pc.onErr != nil {
				pc.onErr(fmt.Errorf("write response error: %w", err))
			}
			if pc.onResRead != nil {
				pc.onResRead(writtenRes)
			}
			return
		}
	}
	if pc.onResRead != nil {
		pc.onResRead(resBytes)
	}
}
