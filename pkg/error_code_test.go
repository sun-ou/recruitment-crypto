package pkg

import (
	"net/http"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestError(t *testing.T) {
	assert.Equal(t, 0, Success.Code(), "Success code should be zero")
	assert.Equal(t, "success", Success.Msg(), "Incrorrect success message")
	assert.Equal(t, "ok", Success.WitchDetails("ok").Details()[0], "Incrorrect success details")

	assert.Equal(t, http.StatusOK, Success.StatusCode(), "Success.StatusCode should be 200")
	assert.Equal(t, http.StatusInternalServerError, ServerError.StatusCode(), "ServerError.StatusCode should be 500")
	assert.Equal(t, http.StatusBadRequest, InvaildParams.StatusCode(), "InvaildParams.StatusCode should be 400")
	assert.Assert(t, strings.Contains(InvaildParams.Error(), "code: 100002"), "InvaildParams.StatusCode should be 400")
}
