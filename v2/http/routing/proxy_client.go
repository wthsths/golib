package gl_routing

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/teris-io/shortid"
)

type responseWriter interface {
	WriteCustomJsonResponse(w http.ResponseWriter, statusCode int, res interface{}) (writtenRes []byte, err error)
}

type ProxyClient struct {
	routeTable       *RouteTable
	routeUrl         string
	httpCli          *http.Client
	shortIDGenerator *shortid.Shortid
	responseWriter   responseWriter

	onErr     func(error, string)
	onReqRead func([]byte, string)
	onResRead func([]byte, string)
}

// NewProxyClient creates a new proxy client instance.
// Underlying HandleRequestAndRedirect method can be registered as a handler function.
// Handler function will redirect incoming request to the routeUrl.
//
// Since it will be registered to a blocking function (ListAndServe), internally occuring event data can be captured via following hooks:
//
// onErr: Can be registered to receive internal errors.
//
// onReqRead: Can be registered to get incoming request body.
//
// onResRead: Can be registered to get outgoing response body.
//
// Hooks will contain a second string value which represents session ID.
// Events which output the same session ID belong to same http session.
//
// RouteTable entries with route parameters can be defined in following fashion:
//
// E.g.: /Transfer/{guid}
func NewProxyClient(routeTable *RouteTable, routeUrl string, httpCli *http.Client, responseWriter responseWriter, onErr func(error, string), onReqRead func([]byte, string), onResRead func([]byte, string)) (*ProxyClient, error) {
	shortIDGenerator, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize short ID generator: %s", err.Error())
	}
	return &ProxyClient{
		routeTable:       routeTable,
		routeUrl:         routeUrl,
		httpCli:          httpCli,
		shortIDGenerator: shortIDGenerator,
		responseWriter:   responseWriter,

		onErr:     onErr,
		onReqRead: onReqRead,
		onResRead: onResRead,
	}, nil
}

// HandleRequestAndRedirect can be registered to http.Handle() for redirecting requests to desired url.
func (pc *ProxyClient) HandleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	sessionID, err := pc.shortIDGenerator.Generate()
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(err, sessionID)
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(err, sessionID)
			}
			return
		}
		if pc.onResRead != nil {
			pc.onReqRead(writtenRes, sessionID)
		}
		return
	}

	uri := r.URL.RequestURI()

	isAllowed := false

	for _, e := range pc.routeTable.routeRules {
		if e.method != r.Method {
			continue
		}

		regexConv, err := RouteToRegExp(uri)
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(err, sessionID)
			}
			writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "internal server error",
			})
			if err != nil {
				if pc.onErr != nil {
					pc.onErr(err, sessionID)
				}
				return
			}
			if pc.onResRead != nil {
				pc.onReqRead(writtenRes, sessionID)
			}
		}

		if e.regexp.MatchString(regexConv) {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("path is not allowed: %s", uri), sessionID)
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"message": "unauthorized call",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(fmt.Errorf("write response error: %s", err.Error()), sessionID)
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes, sessionID)
		}
		return
	}

	redirectUrl := pc.routeUrl + r.URL.RequestURI()
	parsedRedirectUrl, err := url.Parse(redirectUrl)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("unable to parse URL: '%s' error: %s", redirectUrl, err.Error()), sessionID)
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(fmt.Errorf("write response error: %s", err.Error()), sessionID)
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes, sessionID)
		}
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("error reading request body: %s", err.Error()), sessionID)
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(fmt.Errorf("write response error: %s", err.Error()), sessionID)
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes, sessionID)
		}
		return
	}

	if pc.onReqRead != nil {
		pc.onReqRead(reqBytes, sessionID)
	}

	buffer := bytes.NewBuffer(reqBytes)
	nopCloser := ioutil.NopCloser(buffer)

	httpReq := &http.Request{
		Method: r.Method,
		URL:    parsedRedirectUrl,
		Header: r.Header,
		Body:   nopCloser,
	}

	httpRes, err := pc.httpCli.Do(httpReq)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("error executing http request: %s", err.Error()), sessionID)
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(fmt.Errorf("write response error: %s", err.Error()), sessionID)
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes, sessionID)
		}
		return
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("error reading response payload: %s", err.Error()), sessionID)
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(fmt.Errorf("write response error: %s", err.Error()), sessionID)
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes, sessionID)
		}
		return
	}
	defer httpRes.Body.Close()

	for k, v := range httpRes.Header {
		for i := 0; i < len(v); i++ {
			w.Header().Add(k, v[i])
		}
	}

	if httpRes.StatusCode != http.StatusOK {
		w.WriteHeader(httpRes.StatusCode)
	}

	_, err = w.Write(resBytes)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(fmt.Errorf("error writing server response for client: %s", err.Error()), sessionID)
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(fmt.Errorf("write response error: %s", err.Error()), sessionID)
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(writtenRes, sessionID)
		}
		return
	}

	if httpRes.Header.Get("Content-Encoding") == "gzip" {
		reader := bytes.NewReader(resBytes)
		gzipReader, err := gzip.NewReader(reader)
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(fmt.Errorf("error creating gzip reader: %s", err.Error()), sessionID)
			}
			writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "internal error",
			})
			if err != nil {
				if pc.onErr != nil {
					pc.onErr(fmt.Errorf("write response error: %s", err.Error()), sessionID)
				}
				return
			}
			if pc.onResRead != nil {
				pc.onResRead(writtenRes, sessionID)
			}
			return
		}
		/* Modifying resBytes for logging decompressed content AFTER we've written the response body. */
		resBytes, err = ioutil.ReadAll(gzipReader)
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(fmt.Errorf("error reading from gzip reader: %s", err.Error()), sessionID)
			}
			writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "internal error",
			})
			if err != nil {
				if pc.onErr != nil {
					pc.onErr(fmt.Errorf("write response error: %s", err.Error()), sessionID)
				}
				return
			}
			if pc.onResRead != nil {
				pc.onResRead(writtenRes, sessionID)
			}
			return
		}
	}
	if pc.onResRead != nil {
		pc.onResRead(resBytes, sessionID)
	}
}
