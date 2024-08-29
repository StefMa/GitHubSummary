package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/ioki-mobility/summaraizer"
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

type RequestData struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Issue string `json:"issue"`
}

func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	// Read the owner, repo and issue from a POST json that has owner, repo and issue fields
	// in golang.
	// The Request should have everything included we need
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(`{ "summarization": "houston we have a problem" }`))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var requestData RequestData
	err = json.Unmarshal(bodyBytes, &requestData)
	if err != nil {
		w.Write([]byte(`{ "summarization": "houston we have a problem" }`))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestData.Owner != "" && requestData.Repo != "" && requestData.Issue != "" {
		summary, err := fetchAndSummarize(requestData.Owner, requestData.Repo, requestData.Issue)
		if err != nil {
			w.Write([]byte(`{ "summarization": "houston we have a problem" }`))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		summaryBase64 := base64.StdEncoding.EncodeToString([]byte(summary))
		w.Write([]byte(`{ "summarization": "` + summaryBase64 + `" }`))
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Write([]byte(`{ "summarization": "houston we have a problem" }`))
	w.WriteHeader(http.StatusForbidden)
}

func fetchAndSummarize(owner, repo, issue string) (string, error) {
	buffer := bytes.Buffer{}
	gh := summaraizer.GitHub{
		RepoOwner:   owner,
		RepoName:    repo,
		IssueNumber: issue,
	}
	err := gh.Fetch(&buffer)
	if err != nil {
		return "", err
	}

	openAi := summaraizer.OpenAi{
		Model:    "gpt-4o-mini",
		Prompt:   aiPrompt,
		ApiToken: openAiToken,
	}
	summarization, err := openAi.Summarize(&buffer)
	if err != nil {
		return "", err
	}

	return summarization, nil
}
