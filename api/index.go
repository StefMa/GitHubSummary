package api

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*
var templates embed.FS

type templateInfo struct {
	Owner string
	Repo  string
	Issue string
	Error bool
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")
	issue := r.URL.Query().Get("issue")

	tmpl := template.Must(template.ParseFS(templates, "templates/comment.html"))
	var tpl bytes.Buffer
	if owner != "" && repo != "" && issue != "" {
		tmpl.Execute(&tpl, templateInfo{owner, repo, issue, false})
		w.Write(tpl.Bytes())
		w.WriteHeader(http.StatusOK)
		return
	}

	tmpl.Execute(&tpl, templateInfo{owner, repo, issue, true})
	w.Write(tpl.Bytes())
	w.WriteHeader(http.StatusBadRequest)
}
