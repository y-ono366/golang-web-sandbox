package main

import(
	"log"
	"net/http"
	"sync"
	"path/filepath"
	"text/template"
)

type templateHandler struct {
  once     sync.Once
  filename string
  templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("template",t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
  http.Handle("/", &templateHandler{filename: "chat.html"})
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}

