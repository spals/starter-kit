package handler_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"starter-kit/server/config"
	"starter-kit/server/handler"

	"github.com/stretchr/testify/assert"
)

func TestHTTPServerConfigHandler(t *testing.T) {
	handler := handler.NewHTTPServerConfigHandler(&config.HTTPServerConfig{Port: 18080})
	server := httptest.NewServer(handler)
	defer server.Close()

	assert := assert.New(t)
	resp, respErr := http.Get(server.URL)
	if assert.NoError(respErr) {
		assert.Equal(200, resp.StatusCode)
		// FIXME
		// assert.Equal("application/json", resp.Header.Get("Content-Type"))
	}

	body, bodyErr := ioutil.ReadAll(resp.Body)
	if assert.NoError(bodyErr) {
		respConfig := config.HTTPServerConfig{}
		respConfigErr := json.Unmarshal(body, &respConfig)
		if assert.NoError(respConfigErr) {
			assert.Equal(18080, respConfig.Port)
		}
	}
}
