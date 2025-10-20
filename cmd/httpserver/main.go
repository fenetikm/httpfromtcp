package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/fenetikm/httpfromtcp/internal/request"
	"github.com/fenetikm/httpfromtcp/internal/response"
	"github.com/fenetikm/httpfromtcp/internal/server"
)

const port = 42069
const chunkSize = 32

func myHandler(res *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/video" {
		data, err := os.ReadFile("assets/vim.mp4")
		if err != nil {
			log.Fatal(err)
		}
		res.WriteStatusLine(response.StatusCodeOK)
		headers := response.GetDefaultHeaders(len(data))
		headers.Set("Content-Type", "video/mp4")
		res.WriteHeaders(headers)
		res.WriteBody([]byte(data))
		return
	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		url := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
		rget, err := http.Get("https://httpbin.org/" + url)
		if err != nil {
			log.Fatal(err)
		}
		res.WriteStatusLine(response.StatusCodeOK)
		headers := response.GetDefaultHeaders(0)
		headers.Unset("Content-Length")
		headers.Set("Transfer-Encoding", "chunked")
		headers.Set("Trailer", "X-Content-SHA256,X-Content-Length")
		res.WriteHeaders(headers)
		buf := make([]byte, chunkSize)
		fullBody := make([]byte, 0)
		for {
			n, err := rget.Body.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			fullBody = append(fullBody, buf[:n]...)
			res.WriteChunkedBody(buf[:n])
			if n < chunkSize {
				break
			}
		}
		// res.WriteChunkedBodyDone()
		sum := sha256.Sum256(fullBody)
		trailers := headers
		trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", sum))
		trailers.Set("X-Content-Length", strconv.Itoa(len(fullBody)))
		res.WriteTrailers(trailers)
		return
	}

	if req.RequestLine.RequestTarget == "/yourproblem" {
		res.WriteStatusLine(response.StatusCodeBadRequest)
		body := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
		body += "\n"
		headers := response.GetDefaultHeaders(len(body))
		headers.Set("Content-Type", "text/html")
		res.WriteHeaders(headers)
		res.WriteBody([]byte(body))
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		res.WriteStatusLine(response.StatusCodeInternalServerError)
		body := `<html>
<head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
		body += "\n"
		headers := response.GetDefaultHeaders(len(body))
		headers.Set("Content-Type", "text/html")
		res.WriteHeaders(headers)
		res.WriteBody([]byte(body))
		return
	}

	res.WriteStatusLine(response.StatusCodeOK)
	body := `<html>
 <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	body += "\n"
	headers := response.GetDefaultHeaders(len(body))
	headers.Set("Content-Type", "text/html")
	res.WriteHeaders(headers)
	res.WriteBody([]byte(body))
}

func main() {
	server, err := server.Serve(port, myHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
