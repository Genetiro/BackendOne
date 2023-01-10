package transport

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	lr := LinkResources{}

	server := httptest.NewServer(lr.Routes())
	defer server.Close()

	w := httptest.NewRecorder()
	postReq := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader("s=https://www.domain.com/testLink"))
	lr.Create(w, postReq)

}

func TestList(t *testing.T) {
	lr := LinkResources{}

	server := httptest.NewServer(lr.Routes())
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/links", nil)
	w := httptest.NewRecorder()

	lr.List(w, req)

}
