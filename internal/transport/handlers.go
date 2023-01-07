package transport

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/Genetiro/BackendOne/internal/shortner"
	"github.com/go-chi/chi"
)

type ListDb struct {
	Id    int
	link  string
	short string
}

var dbase *sql.DB

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
		r.Get("/", rs.Get)       // GET /links/{short}
		r.Delete("/", rs.Delete) // DELETE /links/{short}

	})
	return r
}

// @Summary List
// @Description get table of links
// @Success     200 {array}   []ListDb{} "list"
// @Failure		400	{string}	string	"ok"
// @Failure		404	{string}	string	"ok"
// @Failure		500	{string}	string	"ok"
// @Router /links [get]
func (rs LinkResources) List(w http.ResponseWriter, r *http.Request) {

	rows, err := dbase.Query("SELECT * FROM links.linkshort")
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
	tmpl, _ := template.ParseFiles("internal/html/list.html")
	tmpl.Execute(w, ListLinks)
}

// @Summary Create
// @Description create new short links
// @Param  result.Link string string true "new link"
// @Success     200 {struct}    Result{} "list"
// @Failure		400	{string}	string	"ok"
// @Failure		404	{string}	string	"ok"
// @Failure		500	{string}	string	"ok"
// @Router /links [post]
func (rs LinkResources) Create(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("html/home.html")
	result := Result{}
	if !shortner.IsValidUrl(r.FormValue("s")) {
		fmt.Println("No valid url")
		result.Status = "Bad format"
		result.Link = ""
	} else {
		result.Link = r.FormValue("s")
		result.Code = shortner.Shorting()

		dbase.Exec("insert into linkshort (link, short) values ($1, $2)", result.Link, result.Code)
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

// @Summary Get
// @Description get link filtered by short
// @Param short path string true "get link"
// @Success     200 {array}   []ListDb{}  "link"
// @Failure		400	{string}	string	"ok"
// @Failure		404	{string}	string	"ok"
// @Failure		500	{string}	string	"ok"
// @Router /links/{short} [get]
func (rs LinkResources) Get(w http.ResponseWriter, r *http.Request) {
	shrt := r.Context().Value("short").(string)

	row := dbase.QueryRow("SELECT id, link, short FROM links.linkshort WHERE short = $1", shrt)

	l := ListDb{}
	errr := row.Scan(&l.Id, &l.link, &l.short)
	switch errr {
	case sql.ErrNoRows:
		fmt.Println("no rows were returned")
		return
	case nil:
		fmt.Println(l)
	default:
		panic(errr)
	}

	tmpl, _ := template.ParseFiles("html/list.html")
	tmpl.Execute(w, l)
}

// @Summary Delete
// @Description delete line from table filtered by short
// @Param short path string true "delete link"
// @Success 301 {integer} integer 1
// @Failure		400	{string}	string	"ok"
// @Failure		404	{string}	string	"ok"
// @Failure		500	{string}	string	"ok"
// @Router /links/{short} [delete]
func (rs LinkResources) Delete(w http.ResponseWriter, r *http.Request) {
	shrt := r.Context().Value("short").(string)

	_, errr := dbase.Exec("DELETe FROM links.linkshort WHERE short = $1", shrt)
	if errr != nil {
		log.Println(errr)
	}
	http.Redirect(w, r, "/", 301)
}
