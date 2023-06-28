package snippet

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetReader(t *testing.T) {
	r, err := getLinkReader("http://vc.ru", 0, nil)
	assert.NoError(t, err)
	defer r.Close()

	body, err := io.ReadAll(r)
	assert.NoError(t, err)
	t.Log(string(body))
}

func TestSnippet(t *testing.T) {
	snippet, err := snippet("http://vc.ru", 0, nil)
	assert.NoError(t, err)
	fmt.Println("Title:", snippet.Title, "\nDescription:", snippet.Description)
}
