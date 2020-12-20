package context

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// HTTPRequestWithContext demonstrates the usage of context in HTTP client request.
func HTTPRequestWithContext(httpPort int, ctxTimeout time.Duration) {
	go func() {
		http.HandleFunc("/slow-query", func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second * 2) // to delay the response
			_, err := w.Write([]byte("hello"))
			if err != nil {
				log.Printf("failed to write response: %v", err)
			}
		})
		if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
			log.Fatalf("failed to start http server: %v", err)
		}
	}()

	time.Sleep(time.Second) // wait for the server to start

	ctx, _ := context.WithTimeout(context.Background(), ctxTimeout) // ctx will Done after ctxTimeout
	url := fmt.Sprintf("http://127.0.0.1:%d/slow-query", httpPort)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatalf("failed to new request: %v", err)
	}

	res, err := http.DefaultClient.Do(req) // will wait for 2 seconds
	if err != nil {
		log.Fatalf("failed to do request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("failed to read body: %v", err)
	}
	fmt.Printf("response: %s\n", string(body))
}
