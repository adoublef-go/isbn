package fetch

import (
	"net/http"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
	// "topheruk.com/encode-json/domain/googleapi" <- use a non imported Go structure
)

func TestFetchUnsplash(t *testing.T) {
	key, path := "edD4pxoguNcEdb_XkdmIOVUBJ9jfPTSLOMAeOPwYZF4", "/photos/random"

	var code int

	opt := DefaultOptions.Header(http.Header{"Accept-Version": []string{"v1"}, "Authorization": []string{`Client-ID ` + key}})

	_ = Fetch("https://api.unsplash.com"+path, func(resp *http.Response) error {
		code = resp.StatusCode
		return nil
	}, opt)

	assert.Assert(t, cmp.Equal(code, 200))
}
