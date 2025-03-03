package server

import (
	"authentification_service/storage"
	"authentification_service/totp"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseCode struct {
	Status bool `json:"status"`
}

type ResponseUser struct {
	Status bool   `json:"status"`
	Code   string `json:"code"`
}

func CheckTOTPCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Login string `json:"login"`
		Code  int    `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	s, err := storage.New()
	if err != nil {
		panic(err)
	}

	code, err := s.GetCode(req.Login)
	if err != nil {
		panic(err)
	}

	if code == "-" {
		w.WriteHeader(http.StatusFailedDependency)
		return
	}

	w.WriteHeader(http.StatusOK)

	is_true_code := totp.CheckTimeBasedCode(req.Login, req.Code)
	resp := ResponseCode{Status: is_true_code}
	b, _ := json.Marshal(resp)

	w.Write(b)
}

func GetUserRandomCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	s, err := storage.New()
	if err != nil {
		panic(err)
	}

	is_correct_user, err := s.CheckUser(req.Login, req.Password)
	if err != nil {
		panic(err)
	}

	if !is_correct_user {
		w.Write([]byte(`"status": false, "message": "There is not user with such data"`))
		return
	}

	code, err := s.GetCode(req.Login)
	if err != nil {
		panic(err)
	}
	if code == "-" {
		w.Write([]byte(`"status": false`))
	} else {
		resp := ResponseUser{Status: true, Code: code}
		b, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		w.Write(b)

	}
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	s, err := storage.New()
	if err != nil {
		panic(err)
	}

	s.SaveUser(req.Login, req.Password)
	w.WriteHeader(http.StatusOK)
}
