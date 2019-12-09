package web

import (
	"encoding/json"
	"net/http"
)

// HTTPStatusCode representa HTTP Status Code
type HTTPStatusCode int

// ErrorMessage represneta mensagem de erro
type ErrorMessage string

func (em ErrorMessage) String() string {
	return string(em)
}

// HTTPResponse representa resposta ao cliente HTTP
type HTTPResponse struct {
	writer  http.ResponseWriter
	code    HTTPStatusCode
	value   interface{}
	message *ErrorMessage
}

// NewHTTPResponse cria instância de HTTPResponse
func NewHTTPResponse(r http.ResponseWriter, c HTTPStatusCode, value interface{}, e *ErrorMessage) (hr *HTTPResponse) {
	hr = new(HTTPResponse)
	hr.writer = r
	hr.code = c
	hr.value = value
	hr.message = e
	return
}

// WriteJSON envia resposta HTTP.
// Converte value em JSON, define status code, e
// enventualmente mensagem de erro.
func (hr *HTTPResponse) WriteJSON() (err error) {

	hr.writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	hr.writer.WriteHeader(int(hr.code))

	if hr.message != nil {
		json.NewEncoder(hr.writer).Encode(
			struct {
				Message string `json:"message"`
			}{
				Message: hr.message.String(),
			},
		)
		return

	}
	json.NewEncoder(hr.writer).Encode(hr.value)

	return
}
