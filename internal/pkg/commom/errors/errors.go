package errors

import (
	"fmt"
)

// Wrap representa comportamento que pode conter erros
type Wrap struct {
	Err error `json:"-"`
}

// Unwrap recupera erro subjacente
func (e *Wrap) Unwrap() error {
	return e.Err
}

// Error representa informações sobre erro
type Error struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// NewError cria instância de Error
func NewError(title, detail string) (e *Error) {
	e = new(Error)
	e.Title = title
	e.Detail = detail
	return
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s[%s]", e.Title, e.Detail)
}

// DomainError representa erros do domínio da aplicação
type DomainError struct {
	Error
}

// NewDomainError cria instância de DomainError
func NewDomainError(title, detail string) (e *DomainError) {
	e = new(DomainError)
	e.Title = title
	e.Detail = detail
	return
}

// ParametersError representa informações sobre erro de parâmetros
type ParametersError struct {
	Title             string           `json:"title"`
	InvalidParameters []ParameterError `json:"invalid-parameters"`
}

// NewParametersError cria instância de ParametersError
func NewParametersError() (e *ParametersError) {
	e = new(ParametersError)
	e.Title = "Um ou Mais parâmetros não são válidos"
	e.InvalidParameters = []ParameterError{}
	return
}

func (e *ParametersError) Error() string {
	return fmt.Sprintf("%s %v", e.Title, e.InvalidParameters)
}

// Add adiciona novo ParameterError
func (e *ParametersError) Add(p ParameterError) {
	e.InvalidParameters = append(e.InvalidParameters, p)
}

// ContainsError informa se existe erros
func (e *ParametersError) ContainsError() bool {
	return len(e.InvalidParameters) > 0
}

// ParameterError representa informações sobre erro de parâmetros
type ParameterError struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Reason string `json:"reason"`
}
