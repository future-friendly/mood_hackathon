package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/future-friednly/mood/backend/models"
	"github.com/future-friednly/mood/backend/util"
	"io/ioutil"
	"net/http"
)

type SignupRequest struct {
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"password"`
}

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := util.DecodeRequest(r.Body, &req); err != nil {
		util.WriteError(w, 500, err)
		return
	}

	if err := models.NewUser(req.Name, req.Email, req.Password); err != nil {
		util.WriteError(w, 500, err)
		return
	}
	util.WriteSuccess(w)
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := util.DecodeRequest(r.Body, &req); err != nil {
		util.WriteError(w, 500, err)
		return
	}

	user, err := models.GetUser(req.Email)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	if err := user.CheckPassword(req.Password); err != nil {
		if errors.Is(err, util.WrongCredentials{}) {
			util.WriteError(w, 401, err)
			return
		}
		util.WriteError(w, 500, err)
		return
	}
	token, err := models.CreateAuthToken(user.ID)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}
	resp, err := json.Marshal(LoginResponse{Token: token})
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	w.WriteHeader(200)
	w.Write(resp)
}

type LogoutRequest struct {
	Token string `json:"token"`
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := util.DecodeRequest(r.Body, &req); err != nil {
		util.WriteError(w, 500, err)
		return
	}

	token, err := models.GetToken(req.Token)
	if err != nil {
		util.WriteError(w, 401, err)
		return
	}
	if err := token.Delete(); err != nil {
		util.WriteError(w, 500, err)
		return
	}
	util.WriteSuccess(w)
}

func includes(value string, arr []string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

var noAuth []string = []string{"/auth/signup", "/auth/login", "/agent/confirm", "/data/newpage"}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if includes(r.RequestURI, noAuth) {
			next.ServeHTTP(w, r)
			return
		}
		req := make(map[string]interface{})
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			util.WriteError(w, 500, err)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if err := json.Unmarshal(body, &req); err != nil {
			util.WriteError(w, 500, err)
			return
		}

		token, ok := req["token"].(string)
		if !ok {
			util.WriteError(w, 401, errors.New("token not present"))
			return
		}
		authToken, err := models.GetToken(token)
		if err != nil {
			util.WriteError(w, 401, err)
			return
		}
		if authToken.IsExpired() {
			util.WriteError(w, 401, errors.New("auth token expired"))
			return
		}
		next.ServeHTTP(w, r)
	})
}