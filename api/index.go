package api

import (
    "bytes"
    "github.com/ioki-mobility/summaraizer"
    "log"
    "net/http"
    "os"
)

const aiPrompt = `
I give you a discussion and you give me a summary.
Each comment of the discussion is wrapped in a <comment author="author-name"> tag.
Your summary should not be longer than 1200 chars.
Here is the discussion:
{{ range $comment := . }}
<comment author="{{ $comment.Author }}">{{ $comment.Body }}</comment>
{{end}}
`

var openAiToken = os.Getenv("OPENAI_API_TOKEN")

func Handler(w http.ResponseWriter, r *http.Request) {
    owner := r.URL.Query().Get("owner")
    repo := r.URL.Query().Get("repo")
    issue := r.URL.Query().Get("issue")
    
    if owner != "" && repo != "" && issue != "" {
        summary := fetchAndSummarize(owner, repo, issue)
        w.Write([]byte(summary))
        w.WriteHeader(http.StatusOK)
        return
    }

    w.Write([]byte("houston we have a problem"))
    w.WriteHeader(http.StatusForbidden)
}

func fetchAndSummarize(owner, repo, issue string) string {
    buffer := bytes.Buffer{}
    gh := summaraizer.GitHub{
        RepoOwner:   owner,
        RepoName:    repo,
        IssueNumber: issue,
    }
    gh.Fetch(&buffer)

    openAi := summaraizer.OpenAi{
        Model:    "gpt-4o-mini",
        Prompt:   aiPrompt,
        ApiToken: openAiToken,
    }
    summarization, err := openAi.Summarize(&buffer)
    if err != nil {
        log.Fatal(err)
    }

    return summarization
}
