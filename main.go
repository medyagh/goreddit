package main

import (
	"log"
	"net/http"
	"html/template"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
)

//basepath for template files
var tmplPath = "./views/templates/"

var userSession *sessions.Session
var MustLogin,_ = template.ParseFiles(tmplPath + "plz-login.html")
func frontPage(w http.ResponseWriter, r *http.Request) {
	data := getTableItems()
	t, _ := template.ParseFiles(tmplPath + "fp.html")
	payload := ScanData{Topics: data, User: userSession.Values["user"]}
	t.Execute(w, payload)
}

func newPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(tmplPath + "mk-topic.html")
	if userSession.Values["user"] != nil {
		payload := User{User: userSession.Values["user"]}
		t.Execute(w, payload)
	}else{
		MustLogin.Execute(w, r)
	}
	
}

func newComment(w http.ResponseWriter, r *http.Request) {
	putComment(r)
}
func showPage(w http.ResponseWriter, r *http.Request) {
	data := getItem(r.URL.Query().Get(":id"))

	payload := GetData{Topic: data, User: userSession.Values["user"],}

	t, _ := template.ParseFiles(tmplPath + "topic.html")
	log.Println(data.Comments)
	t.Execute(w, payload)
}

func upVote(w http.ResponseWriter, r *http.Request) {
	upVoteModel(r)
	http.Redirect(w, r, r.FormValue("returnAddr"), 301)
}

func downVote(w http.ResponseWriter, r *http.Request) {
	downVoteModel(r)
	http.Redirect(w, r, r.FormValue("returnAddr"), 301)
}
func main() {
	fs := http.FileServer(http.Dir("views/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//initialize authentication
	initAuth()

	p := pat.New()

	p.Get("/auth/{provider}/callback", completeAuth)

	p.Get("/auth/{provider}", beginAuth)
	p.Get("/new", newPage)
	p.Post("/create", putItem)
	p.Post("/comment", newComment)
	p.Post("/vote-up", upVote)
	p.Post("/vote-down", downVote)
	p.Get("/topic/id={id}", showPage)
	p.Get("/", frontPage)

	http.Handle("/", p)
	log.Println("Server listening at PORT: 8080")
	http.ListenAndServe(":8080", nil)
}