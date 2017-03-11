package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/enaeseth/asciist/fixture"
)

var s = New()
var contentTypePattern = regexp.MustCompile(`^application/json(; charset=utf-8)?$`)

func TestFixtures(t *testing.T) {
	filenames := []string{"diag-ramp.gif", "bmo.png", "forest.jpg"}

	for _, f := range filenames {
		img, width, art := readFixture(f)
		request := &Request{Image: img, Width: width}

		response := processRequest(request)

		if response.StatusCode != http.StatusOK {
			t.Errorf("%s: got status %s instead of 200", f, response.Status)
		}

		success := Success{}
		if err := readJSON(response, &success); err != nil {
			t.Errorf("%s: %v", f, err)
		}

		if success.Art != art {
			t.Errorf("%s: unexpected art:\n%s", f, success.Art)
		}
	}
}

func TestEmpty(t *testing.T) {
	request := &Request{}
	response := processRequest(request)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("got status %s instead of 400", response.Status)
	}

	failure := Failure{}
	if err := readJSON(response, &failure); err != nil {
		t.Error(err)
	}

	if !strings.Contains(failure.Error, "no image") {
		t.Errorf(`got %#v; expected an Error with "no image"`, failure.Error)
	}
}

func TestNoWidth(t *testing.T) {
	request := &Request{Image: []byte("...")}
	response := processRequest(request)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("got status %s instead of 400", response.Status)
	}

	failure := Failure{}
	if err := readJSON(response, &failure); err != nil {
		t.Error(err)
	}

	if !strings.Contains(failure.Error, "no width") {
		t.Errorf(`got %#v; expected an Error with "no width"`, failure.Error)
	}
}

func TestInvalidImage(t *testing.T) {
	request := &Request{Width: 80, Image: bytes.Repeat([]byte("."), 1024)}
	response := processRequest(request)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("got status %s instead of 400", response.Status)
	}

	failure := Failure{}
	if err := readJSON(response, &failure); err != nil {
		t.Error(err)
	}

	if !strings.Contains(failure.Error, "unknown") {
		t.Errorf(`got %#v; expected an Error with "unknown"`, failure.Error)
	}
}

func processRequest(req *Request) *http.Response {
	payload, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	httpRequest := httptest.NewRequest("POST", "/", bytes.NewReader(payload))
	httpRequest.Header.Set("Content-Type", "application/json")

	return processHTTP(httpRequest)
}

func processHTTP(req *http.Request) *http.Response {
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	return w.Result()
}

func readJSON(response *http.Response, v interface{}) error {
	ct := response.Header.Get("Content-Type")
	if !contentTypePattern.MatchString(ct) {
		return fmt.Errorf("expected JSON but got Content-Type %s", ct)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal %T: %v", v, err)
	}

	return nil
}

func readFixture(imgFilename string) (img []byte, width uint, art string) {
	imgReader, width, art := fixture.ReadFixture(imgFilename)
	img, _ = ioutil.ReadAll(imgReader)
	return
}
