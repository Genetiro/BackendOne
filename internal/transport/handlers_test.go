package transport

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Genetiro/BackendOne/internal/database"
)

func TestCreate(t *testing.T) {
	linksRepo, err := database.NewDB("../database/links.db")
	if err != nil {
		log.Fatal(err)
	}
	defer linksRepo.CloseDb()

	lr := LinkResources{Repo: linksRepo}

	w := httptest.NewRecorder()
	postReq := httptest.NewRequest(http.MethodPost, "/links", strings.NewReader("src_url=https://www.domain.com/testLink"))
	lr.Create(w, postReq)

}

func TestList(t *testing.T) {
	linksRepo, err := database.NewDB("../database/links.db")
	if err != nil {
		log.Fatal(err)
	}
	defer linksRepo.CloseDb()

	lr := LinkResources{Repo: linksRepo}

	req := httptest.NewRequest(http.MethodGet, "/links", nil)
	w := httptest.NewRecorder()

	lr.List(w, req)

}
