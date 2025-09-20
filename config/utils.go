package config

import "net/http"

func DBErrToServiceCode(err error) ServiceCode {
	switch err {
	case nil:
		return CodeOK
	case ErrRecordNotFound:
		return CodeNotFound
	default:
		return CodeDatabaseError
	}
}

func CodeToHttpStatus(code ServiceCode) int {
	switch code {
	case CodeOK:
		return http.StatusOK
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnprocessableEntity, CodeDatabaseError, CodeExternalRequestFail:
		return http.StatusUnprocessableEntity
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeForbidden:
		return http.StatusForbidden
	case CodeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
