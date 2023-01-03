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
		r.Get("/", rs.Get)       // GET /links/{short}
		r.Delete("/", rs.Delete) // DELETE /links/{short}

	})
	return r
}

// @Summary List
// @Description get table of links
// @Param input body ListLinks true "links table"
// @Router /links [get]

func (rs LinkResources) List(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "links.db")
	if err != nil {
		log.Println("can not open db ", err)
	}
	database := db
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
	tmpl, _ := template.ParseFiles("internal/html/list.html")
	tmpl.Execute(w, ListLinks)
}

// @Summary Create
// @Description create new short links
// @Param  input body result true "new short link"
// @Router /links [post]

func (rs LinkResources) Create(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("html/home.html")
	result := Result{}
	if !isValidUrl(r.FormValue("s")) {
		fmt.Println("No valid url")
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

// @Summary Get
// @Description get link filtered by short
// @Param  input body l true "link table"
// @Router /links/{short} [get]

func (rs LinkResources) Get(w http.ResponseWriter, r *http.Request) {
	shrt := r.Context().Value("short").(string)
	db, err := sql.Open("sqlite3", "links.db")
	if err != nil {
		log.Println("can not open db(by short) ", err)
	}
	database := db
	defer db.Close()
	row := database.QueryRow("SELECT id, link, short FROM links.linkshort WHERE short = $1", shrt)

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
// @Param  input body ListLink true
// @Success 301 {integer} integer 1
// @Router /links/{short} [delete]
func (rs LinkResources) Delete(w http.ResponseWriter, r *http.Request) {
	shrt := r.Context().Value("short").(string)
	db, err := sql.Open("sqlite3", "links.db")
	if err != nil {
		log.Println("can not open db(for delete) ", err)
	}
	database := db
	defer db.Close()
	_, errr := database.Exec("DELETe FROM links.linkshort WHERE short = $1", shrt)
	if errr != nil {
		log.Println(errr)
	}
	http.Redirect(w, r, "/", 301)
}
