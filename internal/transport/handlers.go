package transport

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
)

type ListDb struct {
	Id    int
	link  string
	short string
}

var database *sql.DB

type LinkResources struct{}
type Result struct {
	Link   string
	Code   string
	Status string
}

func (rs LinkResources) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", rs.List)    //GET /links
	r.Post("/", rs.Create) //POST /links
	r.Route("/{short}", func(r chi.Router) {
		r.Use(PostCtx)
		r.Get("/", rs.Get)        // GET /links/{short}
		rs.Delete("/", rs.Delete) // DELETE /links/{short}

	})
	return r
}
func (rs LinkResources) List(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "links.db")
	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()
	rows, err := database.Query("SELECT * FROM links.linkshort")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	ListLinks := []ListDb{}
	for rows.Next() {
		l := ListDb{}
		err := rows.Scan(&l.Id, &l.link, &l.short)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ListLinks = append(ListLinks, l)
	}
	tmpl, _ := template.ParseFiles("html/list.html")
	tmpl.Execute(w, ListLinks)
}
func (rs LinkResources) Create(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("html/home.html")
	result := Result{}
	if !isValidUrl(r.FormValue("s")) {
		fmt.Println("Something wrong")
		result.Status = "Bad format"
		result.Link = ""
	} else {
		result.Link = r.FormValue("s")
		result.Code = shorting()
		db, err := sql.Open("sqlite3", "links.linkshort")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		db.Exec("insert into linkshort (link, short) values ($1, $2)", result.Link, result.Code)
		result.Status = "Successfully shorting link"
	}

	tmpl.Execute(w, result)
}
func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "short", chi.URLParam(r, "short"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (rs LinkResources) Get(w http.ResponseWriter, r *http.Request) {

}
func (rs LinkResources) Delete(w http.ResponseWriter, r *http.Request) {

}
