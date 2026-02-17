package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func StreamRouter(li chan *bytes.Buffer) chi.Router {
	r := chi.NewRouter()
	r.Get("/", stream(li))
	return r
}

func stream(li chan *bytes.Buffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		now := time.Now()

		location, err := time.LoadLocation("America/New_York")
		if err != nil {
			fmt.Println("Error loading location:", err)
			return
		}
		file, err := os.OpenFile("clients.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("Failed to open log file:", err)
		}
		defer file.Close()

		log.SetOutput(file)

		<-li

		const boundary = "frame"
		log.Printf("Client connected at %s \n", now.In(location))
		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary)
		w.WriteHeader(http.StatusOK)

		mw := multipart.NewWriter(w)
		mw.SetBoundary(boundary)

		for {
			img := <-li
			data := img.Bytes()

			part, err := mw.CreatePart(textproto.MIMEHeader{
				"Content-Type":   []string{"image/jpeg"},
				"Content-Length": []string{strconv.Itoa(len(data))},
			})
			if err != nil {
				return
			}

			if _, err := part.Write(data); err != nil {
				return
			}
		}
	}
}
