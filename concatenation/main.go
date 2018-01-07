package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func concatenateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	leftPictureURL := r.PostFormValue("left_picture")
	rightPictureURL := r.PostFormValue("right_picture")
	if leftPictureURL == "" || rightPictureURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "left_picture and right_picture are required")
		return
	}
	concatenatedBuffer, err := Concatenate(leftPictureURL, rightPictureURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	_, err = concatenatedBuffer.WriteTo(w)
	if err != nil {
		log.Errorf("Can't return the image data: %s", err)
	}
}

func main() {
	r := httprouter.New()
	r.POST("/concatenate", concatenateHandler)
	log.Info("start listening")
	err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), r)
	if err != nil {
		panic(errors.Errorf("Cannot run server: %s", err))
	}
}
