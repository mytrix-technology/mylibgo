package fasthttp

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"os"
	"reflect"
	"time"
)

var (
	headerContentTypeJson = []byte("application/json")
	client                *fasthttp.Client
	domains               = make(map[string]fasthttp.RequestHandler)
)

type Entity struct {
	Id   int
	Name string
}

func sendGetRequest(UrlGet string) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(UrlGet)
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err == nil {
		fmt.Printf("DEBUG Response: %s\n", resp.Body())
	} else {
		fmt.Fprintf(os.Stderr, "ERR Connection error: %v\n", err)
	}
	fasthttp.ReleaseResponse(resp)
}

func sendPostRequest(UrlPost string) {
	// per-request timeout
	reqTimeout := time.Duration(100) * time.Millisecond

	reqEntity := &Entity{
		Name: "New entity",
	}
	reqEntityBytes, _ := json.Marshal(reqEntity)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(UrlPost)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.SetBodyRaw(reqEntityBytes)
	resp := fasthttp.AcquireResponse()
	err := client.DoTimeout(req, resp, reqTimeout)
	fasthttp.ReleaseRequest(req)
	if err == nil {
		statusCode := resp.StatusCode()
		respBody := resp.Body()
		fmt.Printf("DEBUG Response: %s\n", respBody)
		if statusCode == http.StatusOK {
			respEntity := &Entity{}
			err = json.Unmarshal(respBody, respEntity)
			if err == io.EOF || err == nil {
				fmt.Printf("DEBUG Parsed Response: %v\n", respEntity)
			} else {
				fmt.Fprintf(os.Stderr, "ERR failed to parse response: %v\n", err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "ERR invalid HTTP response code: %d\n", statusCode)
		}
	} else {
		errName, known := httpConnError(err)
		if known {
			fmt.Fprintf(os.Stderr, "WARN conn error: %v\n", errName)
		} else {
			fmt.Fprintf(os.Stderr, "ERR conn failure: %v %v\n", errName, err)
		}
	}
	fasthttp.ReleaseResponse(resp)
}

func httpConnError(err error) (string, bool) {
	errName := ""
	known := false
	if err == fasthttp.ErrTimeout {
		errName = "timeout"
		known = true
	} else if err == fasthttp.ErrNoFreeConns {
		errName = "conn_limit"
		known = true
	} else if err == fasthttp.ErrConnectionClosed {
		errName = "conn_close"
		known = true
	} else {
		errName = reflect.TypeOf(err).String()
		if errName == "*net.OpError" {
			// Write and Read errors are not so often and in fact they just mean timeout problems
			errName = "timeout"
			known = true
		}
	}
	return errName, known
}
