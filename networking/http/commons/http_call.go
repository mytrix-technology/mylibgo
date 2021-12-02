package commons

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type HttpCall struct {
	Method			HttpCallMethod
	URL				string
	DataRequest		[]byte
	ContentType		HttpCallContentType
	Headers			map[string]string
}

type HttpCallContentType string
const (
	HTTP_CALL_CONTENT_XML	HttpCallContentType = "application/xml; charset=utf-8"
	HTTP_CALL_CONTENT_JSON	HttpCallContentType = "application/json"
	HTTP_CALL_CONTENT_FORM	HttpCallContentType = "application/x-www-form-urlencoded"
)

type HttpCallMethod string
const (
	HTTP_CALL_METHOD_GET	HttpCallMethod = "GET"
	HTTP_CALL_METHOD_POST	HttpCallMethod = "POST"
)

func (httpCall *HttpCall) SendRequest() ([]byte, error) {
	httpRequest, err := http.NewRequest(string(httpCall.Method), httpCall.URL, bytes.NewBuffer(httpCall.DataRequest))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Add("Content-Type", string(httpCall.ContentType))
	if httpCall.Headers != nil {
		for headerName, headerValue := range httpCall.Headers {
			httpRequest.Header.Add(headerName, headerValue)
		}
	}

	client := &http.Client{}
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	response, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}
