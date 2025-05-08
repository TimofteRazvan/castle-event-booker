package render

import (
	"net/http"
	"testing"

	"github.com/TimofteRazvan/castle-event-booker/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	request, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(request.Context(), "flash", "123")
	result := AddDefaultData(&td, request)
	if result.Flash != "123" {
		t.Error("Wrong Flash value, expected 123")
	}
}

func getSession() (*http.Request, error) {
	request, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		return nil, err
	}

	context := request.Context()
	context, _ = session.Load(context, request.Header.Get("X-Session"))
	request = request.WithContext(context)

	return request, nil
}
