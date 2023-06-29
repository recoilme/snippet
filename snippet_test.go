package snippet

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetReader(t *testing.T) {
	r, statusCode, err := GetLinkReader("http://vc.ru", 0, nil)
	assert.NoError(t, err)
	defer r.Close()

	body, err := io.ReadAll(r)
	assert.NoError(t, err)
	t.Log(statusCode, string(body))
}

func TestSnippet(t *testing.T) {
	snippet, err := Snippet("http://vc.ru", 0, nil)
	assert.NoError(t, err)
	t.Log("Title:", snippet.Title, "\nDescription:", snippet.Description, "\nstatusCode:", snippet.StatusCode)
}
