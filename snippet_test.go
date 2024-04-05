package snippet

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetReader(t *testing.T) {
	r, statusCode, err := GetLinkReader("http://vc.ru", 0, nil, 0)
	assert.NoError(t, err)
	defer r.Close()

	body, err := io.ReadAll(r)
	assert.NoError(t, err)
	t.Log(statusCode, string(body))
}

func TestSnippet(t *testing.T) {
	snippet, err := Snippet("https://vc.ru/", 0, nil)
	assert.NoError(t, err)
	t.Log("Title:", snippet.Title, "\nDescription:", snippet.Description, "\nstatusCode:", snippet.StatusCode)
}

func TestSnippetMeta(t *testing.T) {
	snippet, err := Snippet("https://vc.ru/hr/1110372", 0, nil)
	assert.NoError(t, err)
	t.Log("Title:", snippet.Title, "\nDescription:", snippet.Description, "\nSection:", snippet.Section, "\nTag:", snippet.Tag, "\nKeywords:", snippet.Keywords, "\nstatusCode:", snippet.StatusCode)
}

func TestSnippetMetaKeywords(t *testing.T) {
	snippet, err := Snippet("https://habr.com/ru/news/805597/", 0, nil)
	assert.NoError(t, err)
	t.Log("Title:", snippet.Title, "\nDescription:", snippet.Description, "\nSection:", snippet.Section, "\nTag:", snippet.Tag, "\nKeywords:", snippet.Keywords, "\nstatusCode:", snippet.StatusCode)
}
