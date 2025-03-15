package totp

import (
	"authentification_service/storage"
)

func CheckTimeBasedCode(login string, send_code int) bool {
	s, err := storage.New()
	if err != nil {
		panic(err)
	}

	_, code, _ := s.GetCode(login)
	true_code := CreateTimeBasedCode(code)
	return true_code == send_code
}
