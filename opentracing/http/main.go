package main

import (
	"context"
	"errors"
	"examples/opentracing/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"io/ioutil"
	"log"
	"net/http"
)

// configurations
const (
	serverHost = "localhost"
	serverPort = 6666

	tracerAgentHost   = "localhost"
	tracerAgentPort   = 32772
	tracerServiceName = "myHTTPService"
)

var httpHost = fmt.Sprintf("http://%s:%d", serverHost, serverPort)

// traceMiddleware retrieve spanContext from client request, if there is, start a span childOf it,
// otherwise, start a new root span.
func traceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if opentracing.IsGlobalTracerRegistered() {
			tracer := opentracing.GlobalTracer()
			clientSpanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
			var (
				span opentracing.Span
				// note: the span operation name should be a general string to allow the tracing systems to do aggregations,
				// so the request path should be a gin registered route path rather than a request URL, since a registered
				// route path may contain path params like '/user/:id'.
				spanOpName = fmt.Sprintf("%s %s", c.Request.Method, c.FullPath() /* not c.Request.URL */)
			)
			if err == nil { // has client span context, inherit it
				span = tracer.StartSpan(spanOpName, opentracing.ChildOf(clientSpanCtx))
			} else { // doesn't has, start a root span
				span = tracer.StartSpan(spanOpName)
			}
			defer span.Finish()

			// set span tag
			ext.Component.Set(span, "gin")
			ext.HTTPMethod.Set(span, c.Request.Method)
			ext.HTTPUrl.Set(span, c.Request.URL.String())
			span.SetTag("http.params", fmt.Sprint(c.Params))
			span.SetTag("http.queries", fmt.Sprint(c.Request.URL.Query()))

			// put the span in request-scope context of gin
			c.Set("span", span)
		}

		c.Next()
	}
}

// HTTPGetWithTrace inject this service scoped spanContext into carrier of outgoing HTTP request,
// to allow the server to extract it from incoming request.
func HTTPGetWithTrace(URL string, c *gin.Context) (body []byte, err error) {
	if !opentracing.IsGlobalTracerRegistered() {
		err = errors.New("HTTPGetWithTrace: global tracer is not registered")
		return
	}
	tracer := opentracing.GlobalTracer()

	v, exist := c.Get("span")
	span, ok := v.(opentracing.Span)
	if !exist || !ok {
		err = errors.New("HTTPGetWithTrace: can't retrieve span from context")
		return
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, URL, nil)
	if err != nil {
		err = fmt.Errorf("HTTPGetWithTrace: can't new request: %w", err)
		return
	}

	if err = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
		err = fmt.Errorf("HTTPGetWithTrace: can't inject spanContext to carrier")
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("HTTPGetWithTrace: can't do request: %w", err)
		return
	}
	defer res.Body.Close()

	if body, err = ioutil.ReadAll(res.Body); err != nil {
		err = fmt.Errorf("HTTPGetWithTrace: can't read body: %w", err)
		return
	}

	return
}

// handlers

func fooHandler(c *gin.Context) {
	barBody, err := HTTPGetWithTrace(httpHost+"/bar", c)
	if err != nil {
		c.String(http.StatusInternalServerError, "%d", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	bazBody, err := HTTPGetWithTrace(httpHost+"/baz", c)
	if err != nil {
		c.String(http.StatusInternalServerError, "%d", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	c.String(http.StatusOK, "%s\n%s\n", barBody, bazBody)
}

func barHandler(c *gin.Context) {
	quxBody, err := HTTPGetWithTrace(httpHost+"/qux", c)
	if err != nil {
		c.String(http.StatusInternalServerError, "%d", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	c.String(http.StatusOK, string(quxBody))
}

func bazHandler(c *gin.Context) {
	c.String(http.StatusOK, "response from baz")
}

func quzHandler(c *gin.Context) {
	c.String(http.StatusOK, "response from qux")
}

func main() {
	closer := utils.InitTracer(tracerServiceName, tracerAgentHost, tracerAgentPort)
	defer closer.Close()

	r := gin.Default()
	r.Use(traceMiddleware()) // use trace

	// network tracing:
	//     /foo
	//       |
	//    |------|
	//   /bar  /baz
	//    |
	//   /qux

	r.GET("/foo", fooHandler)
	r.GET("/bar", barHandler)
	r.GET("/baz", bazHandler)
	r.GET("/qux", quzHandler)

	if err := r.Run(":6666"); err != nil {
		err = fmt.Errorf("main: run gin server: %w", err)
		log.Fatal(err)
	}
}
