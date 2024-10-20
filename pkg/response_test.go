package pkg

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)


func TestToResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w1 := httptest.NewRecorder()
	ctx1, _ := gin.CreateTestContext(w1)

	NewResponse(ctx1).ToResponse(nil)
	if w1.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w1.Code)
	}
	if w1.Body.String() != "{}" {
		t.Errorf("Expected response body '{}', got '%s'", w1.Body.String())
	}

	w2 := httptest.NewRecorder()
	ctx2, _ := gin.CreateTestContext(w2)
	NewResponse(ctx2).ToResponse(map[string]string{"hello": "world"})
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w2.Code)
	}
	if w2.Body.String() != `{"hello":"world"}` {
		t.Errorf(`Expected response body '{"hello":"world"}', got '%s'`, w2.Body.String())
	}
}

func TestToErrorResponseResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	NewResponse(ctx).ToErrorResponse(InvaildParams.WitchDetails("missing parameter user name"))
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	if !strings.Contains(w.Body.String(), `"code":100002`) || !strings.Contains(w.Body.String(), "missing parameter user name") {
		t.Errorf("Unexpected response body, got '%s'", w.Body.String())
	}
}
