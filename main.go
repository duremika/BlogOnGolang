package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"./models"
	"os"
)

var port = ":80"

var posts = make(map[string]*models.Post, 0)

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles(
		"templates/index.html",
		"templates/header.html",
		"templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "index", posts)
}

func writeHandler(w http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles(
		"templates/write.html",
		"templates/header.html",
		"templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "write", nil)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"templates/write.html",
		"templates/header.html",
		"templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	id := r.FormValue("id")
	post, exist := posts[id]
	if !exist{
		http.NotFound(w, r)
	}

	t.ExecuteTemplate(w, "write", post)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == ""{
		http.NotFound(w, r)
	}
	delete(posts,id)

	http.Redirect(w, r, "/", 302)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	var post *models.Post
	if id != ""{
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		for exist := true; exist; _, exist = posts[id] {
			id = GenerateId()
			post = models.NewPost(id, title, content)
		}
		posts[id] = post
	}

	http.Redirect(w, r, "/", 302)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/SavePost", savePostHandler)

}
