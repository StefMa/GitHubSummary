package api

import (
	"bytes"
	"html/template"
	"net/http"
)

func ForOForHandler(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")
	issue := r.URL.Query().Get("issue")

	tmpl := template.Must(template.ParseFS(templates, "templates/404.html"))
	var tpl bytes.Buffer
	tmpl.Execute(&tpl, templateInfo{owner, repo, issue, true})
	w.Write(tpl.Bytes())
	w.WriteHeader(http.StatusBadRequest)
}
