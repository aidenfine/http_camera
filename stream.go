package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func StreamRouter(li chan *bytes.Buffer) chi.Router {
	r := chi.NewRouter()
	r.Get("/", stream(li))
	return r
}

func stream(li chan *bytes.Buffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		<-li

		const boundary = "frame"
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
