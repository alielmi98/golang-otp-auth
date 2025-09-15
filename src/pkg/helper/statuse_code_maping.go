package helper

import (
	"net/http"

	"github.com/alielmi98/golang-otp-auth/pkg/service_errors"
)

var StatusCodeMapping = map[string]int{
	// User
	service_errors.EmailExists:               409,
	service_errors.UsernameExists:            409,
	service_errors.RecordNotFound:            404,
	service_errors.PermissionDenied:          403,
	service_errors.UsernameOrPasswordInvalid: 401,
	// Token
	service_errors.InvalidRefreshToken: 401,
	service_errors.TokenRequired:       401,
	service_errors.TokenExpired:        401,
	service_errors.TokenInvalid:        401,
	service_errors.ClaimsNotFound:      401,
	service_errors.InvalidRolesFormat:  400,
	// OTP
	service_errors.OptExists:   409,
	service_errors.OtpUsed:     400,
	service_errors.OtpNotValid: 400,
	// Validation
	service_errors.ValidationError: 400,
}

func TranslateErrorToStatusCode(err error) int {
	value, ok := StatusCodeMapping[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}
	return value
}
