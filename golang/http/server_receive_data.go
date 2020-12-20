package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var serveMux *http.ServeMux

// RunServerReceiveQsExample demonstrates how a HTTP server receives query string data.
// Take a try after the service starts:
//    GET /test/qs?age=23&name=xavier HTTP/1.1
func RunServerReceiveQsExample(addr string) {
	serveMux.HandleFunc("/test/qs", func(writer http.ResponseWriter, request *http.Request) {
		// the Form and PostForm field of Request is only available after ParseForm called
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, "internal error", http.StatusInternalServerError)
			log.Printf("failed to parseForm: %v", err)
			return
		}

		// the Form field of Request contains not only x-www-form-urlencoded data but also query string data
		name := request.Form.Get("name")
		age := request.Form.Get("age")
		msg := fmt.Sprintf("qs test: %s is %s years old", name, age)
		_, err = writer.Write([]byte(msg))
		if err != nil {
			log.Printf("failed to write: %v", err)
			return
		}
	})

	err := http.ListenAndServe(addr, serveMux)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
}

// RunServerReceiveFormExample demonstrate how a HTTP server receives "application/x-www-form-urlencoded" data.
// Take a try after the service starts:
//    POST /test/form HTTP/1.1
//    Content-Type: application/x-www-form-urlencoded
//    Content-Length: 18
//
//    name=xavier&age=23
func RunServerReceiveFormExample(addr string) {
	serveMux.HandleFunc("/test/form", func(writer http.ResponseWriter, request *http.Request) {
		// the Form and PostForm field of Request is only available after ParseForm called
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, "internal error", http.StatusInternalServerError)
			log.Printf("failed to parseForm: %v", err)
			return
		}

		// the PostForm field of Request contains only x-www-form-urlencoded data
		name := request.PostForm.Get("name")
		age := request.PostForm.Get("age")
		msg := fmt.Sprintf("form test: %s is %s years old", name, age)
		_, err = writer.Write([]byte(msg))
		if err != nil {
			log.Printf("failed to write: %v", err)
			return
		}
	})

	err := http.ListenAndServe(addr, serveMux)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
}

// RunServerReceiveJsonExample demonstrates how a HTTP server receives "application/json" data.
// Take a try after the service starts:
//    POST /test/json HTTP/1.1
//    Content-Type: application/json
//    Content-Length: 29
//
//    {"name": "Saudade","age": 23}
func RunServerReceiveJsonExample(addr string) {
	serveMux.HandleFunc("/test/json", func(writer http.ResponseWriter, request *http.Request) {
		jsonBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "server error", http.StatusInternalServerError)
			log.Printf("failed to read body: %v", err)
			return
		}
		fmt.Println(string(jsonBytes))
		writer.Header().Add("content-type", "application/json; charset=utf-8")
		_, _ = writer.Write(jsonBytes)

		// next, unmarshal the jsonBytes
		// todo ...
	})

	err := http.ListenAndServe(addr, serveMux)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
}

// RunServerReceiveMultiExample demonstrates how a HTTP server receives "multipart/form-data" data (uploaded files).
// Take a try (send a field "foo=bar" and a file) after the service starts:
//    POST /test/multi HTTP/1.1
//    Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW
//    Content-Length: 56442
//
//
//    Content-Disposition: form-data; name="foo"
//
//    bar
//    ------WebKitFormBoundary7MA4YWxkTrZu0gW--,
//    Content-Disposition: form-data; name="foo"
//
//    bar
//    ------WebKitFormBoundary7MA4YWxkTrZu0gW--
//    Content-Disposition: form-data; name="avatarFile"; filename="/Users/xvrzhao/Pictures/avatars/ava.jpg
//
//
//    ------WebKitFormBoundary7MA4YWxkTrZu0gW--
func RunServerReceiveMultiExample(addr string) {
	serveMux.HandleFunc("/test/multi", func(writer http.ResponseWriter, request *http.Request) {
		// request.MultipartForm is available only when ParseMultipartForm has been called
		err := request.ParseMultipartForm(2 * 1024)
		if err != nil {
			http.Error(writer, "server error", http.StatusInternalServerError)
			log.Printf("failed to parse nultipart form: %v", err)
			return
		}

		// Value is type map[string][]string
		foo := request.MultipartForm.Value["foo"][0]

		// File is type map[string][]*multipart.FileHeader
		fileHeader := request.MultipartForm.File["avatarFile"][0]
		srcFile, err := fileHeader.Open()
		if err != nil {
			http.Error(writer, "server error", http.StatusInternalServerError)
			log.Printf("failed to open file: %v", err)
			return
		}
		defer srcFile.Close()

		dstFile, err := os.Create("./" + fileHeader.Filename) // if filename is relative path, then that path is relative to pwd in which exec the bin
		if err != nil {
			http.Error(writer, "server error", http.StatusInternalServerError)
			log.Printf("failed to create dst file: %v", err)
			return
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			http.Error(writer, "server error", http.StatusInternalServerError)
			log.Printf("failed to copy file: %v", err)
			return
		}

		msg := fmt.Sprintf("avatarFile is uploaded, foo received: %s", foo)
		_, _ = writer.Write([]byte(msg))
	})

	err := http.ListenAndServe(addr, serveMux)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
}

func init() {
	serveMux = http.NewServeMux()
}
