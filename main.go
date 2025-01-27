package main

import (
	"authentification_service/storage"
	"fmt"
)

func main() {
	s, _ := storage.New()
	new_user := storage.User{Login: "new_login1", Password: "new_password1"}
	err := s.SaveUser(&new_user)
	if err != nil {
		panic(err)
	}
	code, _ := s.GetCode(&storage.User{Login: "new_login1", Password: "new_password1"})
	fmt.Println(code)
}
