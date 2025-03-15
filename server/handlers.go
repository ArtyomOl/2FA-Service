package server

import (
	"authentification_service/storage"
	"authentification_service/totp"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ResponseCode struct {
	Status  bool   `json:"status"`
	Message string `json:"message"]`
}

type ResponseUser struct {
	Status bool   `json:"status"`
	Code   string `json:"code"`
}

func CheckTOTPCodeHandler(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		defer s.Close()

		status, _, err := s.GetCode(req.Login)
		if err != nil {
			panic(err)
		}

		if !status {
			w.WriteHeader(http.StatusFailedDependency)
			return
		}

		w.WriteHeader(http.StatusOK)

		is_true_code := totp.CheckTimeBasedCode(req.Login, req.Code)

		if !is_true_code {
			val, err := client.Get(context.Background(), req.Login).Result()
			if err == redis.Nil {
				fmt.Println("value not found")
			}
			intval, _ := strconv.Atoi(val)

			if intval > 3 {
				resp := ResponseCode{Status: false, Message: "Too many failed attempts"}
				b, _ := json.Marshal(resp)

				w.Write(b)
				return
			}

			val = strconv.Itoa(intval + 1)
			if err := client.Set(context.Background(), req.Login, val, 30*time.Second).Err(); err != nil {
				fmt.Printf("failed to set data, error: %s", err.Error())
			}
		}

		resp := ResponseCode{Status: is_true_code}
		b, _ := json.Marshal(resp)

		w.Write(b)
	}
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

	status, code, err := s.GetCode(req.Login)
	if err != nil {
		panic(err)
	}
	if !status {
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
	fmt.Println("work!!!!!!!!!!!!!!1")

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
