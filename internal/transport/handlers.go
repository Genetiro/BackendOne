package transport

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Genetiro/BackendOne/internal/database"
	"github.com/Genetiro/BackendOne/internal/shortner"
	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

type LinkResources struct {
	Repo *database.Db
}

func (rs LinkResources) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", rs.List)    //GET /links
	r.Post("/", rs.Create) //POST /links
	//r.Route("/{short}", func(r chi.Router) {
	//r.Use(PostCtx)
	r.Get("/{short}", rs.Get)       // GET /links/{short}
	r.Delete("/{short}", rs.Delete) // DELETE /links/{short}

	//})
	return r
}

// @Summary List
// @Description get table of links
// @Success     200 {array}   array "list"
// @Failure		400	{string}	string	"ok"
// @Failure		404	{string}	string	"ok"
// @Failure		500	{string}	string	"ok"
// @Router /links [get]
func (rs LinkResources) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List")

	list, err := rs.Repo.GetShortLinks()
	if err != nil {
		log.Println("can't get links ", err)
	}
	tmpl, _ := template.ParseFiles("../html/list.html")
	tmpl.Execute(w, list)

}

// @Summary Create
// @Description create new short links
// @Param  result.Link body string true "new link"
// @Success     200 {struct}    Result{} "list"
// @Failure		400	{string}	string	"ok"
// @Failure		404	{string}	string	"ok"
// @Failure		500	{string}	string	"ok"
// @Router /links [post]
func (rs LinkResources) Create(w http.ResponseWriter, r *http.Request) {
	NewLink := database.ListDb{}
	tmpl, err := template.ParseFiles("../html/home.html")
	if err != nil {
		log.Println("can't parse template")
	}
	if !shortner.IsValidUrl(r.FormValue("src_url")) {
		log.Println("No valid url")

	} else {
		NewLink.Link = r.FormValue("src_url")
		NewLink.Short = shortner.Shorting()
		rs.Repo.CreateShort(NewLink)

		// newBdRow, err := rs.Repo.CreateShort(NewLink)
		// if err != nil {
		// 	log.Fatal(err)
	}

	tmpl.Execute(w, NewLink)
}

// func PostCtx(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := context.WithValue(r.Context(), "short", chi.URLParam(r, "short"))
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// @Summary Get
// @Description get link filtered by short
// @Param short path string true "get link"
// @Success     200 {array}   array  "link"
// @Failure		400	{string}	string	"ok"
// @Failure		404	{string}	string	"ok"
// @Failure		500	{string}	string	"ok"
// @Router /links/{short} [get]
func (rs LinkResources) Get(w http.ResponseWriter, r *http.Request) {
	//shrt := r.Context().Value("short").(string)
	fmt.Println("List by short link")

	shrt := chi.URLParam(r, "short")

	l, err := rs.Repo.GetByShort(shrt)
	if err != nil {
		log.Println("have not short link ", err)
	}

	tmpl, _ := template.ParseFiles("../html/list.html")
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
	//shrt := r.Context().Value("short").(string)
	shrt := chi.URLParam(r, "short")

	err := rs.Repo.DeleteShort(shrt)
	if err != nil {
		log.Println("can't delete", err)
	}
	http.Redirect(w, r, "/", 301)
}
