package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createRequest(ctx context.Context, event events.APIGatewayProxyRequest) *http.Request {
	request := &http.Request{
		Method:        event.HTTPMethod,
		Proto:         "HTTP/1.0", // Lambda does not support keep alive neither http2
		ProtoMajor:    1,
		ProtoMinor:    0,
		ContentLength: int64(len([]byte(event.Body))),
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Transfer-Encoding
		TransferEncoding: []string{"identity"},
		Header:           make(http.Header),
		Close:            false,
		Host:             event.Headers["Host"],
		Trailer:          make(http.Header), // Not supported
		RemoteAddr:       event.RequestContext.Identity.SourceIP,
		RequestURI:       event.Path,
		TLS:              nil, // Not supported
	}
	request = request.WithContext(ctx)

	for name, value := range event.Headers {
		request.Header.Add(name, value)
	}

	return request
}

func proxy(ctx context.Context, h http.Handler, event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	request := createRequest(ctx, event)
	responseWriter := httptest.NewRecorder()

	h.ServeHTTP(responseWriter, request)

	return createResponse(ctx, responseWriter)
}

func createResponse(_ context.Context, w *httptest.ResponseRecorder) events.APIGatewayProxyResponse {
	response := events.APIGatewayProxyResponse{
		StatusCode:      w.Code,
		IsBase64Encoded: false,
		Headers:         make(map[string]string),
	}

	for header, entries := range w.HeaderMap {
		response.Headers[header] = strings.Join(entries, ",")
	}

	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		panic(fmt.Sprintf("Body should not be null err=%s", err))
	}

	response.Body = string(data)

	return response
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-MyCustom-Header", "This is the value of my custom header")
	fmt.Fprint(w, "Hello World")
}

func LambdaHandler(h http.Handler) interface{} {
	return func(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		response := proxy(ctx, h, r)
		return response, nil
	}
}

func main() {
	lambda.Start(LambdaHandler(http.HandlerFunc(HelloWorldHandler)))
}
