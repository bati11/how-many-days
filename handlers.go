package how_many_days

import (
	"net/http"
	"time"
	"html/template"
	"github.com/golang/appengine"
	"github.com/golang/appengine/datastore"
	"github.com/satori/go.uuid"
	"strings"
)

type Anniversary struct {
	Date  time.Time
	Title string
}

var inputPageTemplate = template.Must(template.ParseFiles("templates/input.html"))
var showPageTemplate  = template.Must(template.ParseFiles("templates/show.html"))

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/add", add)
	http.HandleFunc("/show/", show)
}

func root(w http.ResponseWriter, r *http.Request) {
	if err := inputPageTemplate.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func add(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "title should be present", http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, "date format should be '2016-05-12'", http.StatusBadRequest)
		return
	}
	anniversary := Anniversary{
		Date: date,
		Title: title,
	}
	stringID := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	key := datastore.NewKey(c, "Anniversary", stringID, 0, nil)
	key, err = datastore.Put(c, key, &anniversary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/show/" + stringID, http.StatusFound)
}

func show(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	uuid := r.URL.Path[6:]
	key := datastore.NewKey(c, "Anniversary", uuid, 0, nil)
	var anniversary Anniversary
	if err := datastore.Get(c, key, &anniversary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := showPageTemplate.Execute(w, &anniversary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
