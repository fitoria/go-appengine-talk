package imgcurl

import (
	"fmt"
	"net/http"

	"appengine"
	"appengine/datastore"

	"github.com/julienschmidt/httprouter"
)

type Link struct {
	Url   string
	Alias string
}

func LinkKey(c appengine.Context, alias string) *datastore.Key {
	return datastore.NewKey(c, "Link", alias, 0, nil)
}

func init() {
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/", postLink)
	router.GET("/:alias", getLink)
	http.Handle("/", router)
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Imgcurl app engine demo")
}

func postLink(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	link := Link{
		Url:   r.FormValue("url"),
		Alias: r.FormValue("alias"),
	}

	c := appengine.NewContext(r)
	key := LinkKey(c, link.Alias)
	_, err := datastore.Put(c, key, &link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprint(w, "Link Added")
	}
}

func getLink(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	alias := params.ByName("alias")
	c := appengine.NewContext(r)

	key := LinkKey(c, alias)
	fmt.Print(key)
	link := new(Link)
	if err := datastore.Get(c, key, link); err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	http.Redirect(w, r, link.Url, http.StatusFound)
}
