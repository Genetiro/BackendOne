package transport

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Genetiro/BackendOne/internal/database"
	"github.com/Genetiro/BackendOne/internal/shortner"
	"github.com/go-chi/chi"
)

type LinkResources struct{}

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
	base, err := database.NewDB("./links.db")
	if err != nil {
		log.Println("can't connect for get list ", err)
	}
	list, _ := base.GetShortLinks()
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
	if !shortner.IsValidUrl(r.FormValue("s")) {
		log.Println("No valid url")

	} else {
		NewLink.Link = r.FormValue("s")
		NewLink.Short = shortner.Shorting()

	}
	base, err := database.NewDB("./links.db")
	if err != nil {
		log.Println("can't connect for get list ", err)
	}
	newBdRow, err := base.CreateShort(NewLink)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, newBdRow)
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
	shrt := chi.URLParam(r, "short")
	base, err := database.NewDB("./links.db")
	if err != nil {
		log.Println("can't connect for get list ", err)
	}
	l, errr := base.GetByShort(shrt)
	if errr != nil {
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
	base, err := database.NewDB("./links.db")
	if err != nil {
		log.Println("can't connect for get list ", err)
	}
	errr := base.DeleteShort(shrt)
	if errr != nil {
		log.Println("can't delete", errr)
	}
	http.Redirect(w, r, "/", 301)
}
