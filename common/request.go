package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// DrainBody drains the body of the given of the given HTTP request.
func DrainBody(r *http.Request) {

	if r.Body == nil {
		return
	}

	_, err := io.Copy(ioutil.Discard, r.Body)
	if err != nil {
		fmt.Printf("error draining request body: [%v]", err)
	}

	err = r.Body.Close()
	if err != nil {
		fmt.Printf("error closing request body: [%v]", err)
	}
}
