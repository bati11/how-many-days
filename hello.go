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
	Uuid  string
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
	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	anniversary := Anniversary{
		Uuid: strings.Replace(uuid.NewV4().String(), "-", "", -1),
		Date: date,
		Title: r.FormValue("title"),
	}
	key := datastore.NewIncompleteKey(c, "Anniversary", nil)
	key, err = datastore.Put(c, key, &anniversary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/show/" + anniversary.Uuid, http.StatusFound)
}

func show(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	uuid := r.URL.Path[6:]
	var annversaries []Anniversary
	q := datastore.NewQuery("Anniversary").Filter("Uuid =", uuid)
	if _, err := q.GetAll(c, &annversaries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(annversaries) == 0 {
		http.NotFound(w, r);
		return
	}
	annversary := annversaries[0]
	if err := showPageTemplate.Execute(w, &annversary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}