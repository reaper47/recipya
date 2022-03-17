package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/config"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/regex"
	"github.com/reaper47/recipya/internal/repository"
	"github.com/reaper47/recipya/internal/templates"
)

// Register handles the [GET,POST] /auth/register endpoint.
func Register(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		handleGetRegister(w, req)
	case http.MethodPost:
		handlePostRegister(w, req)
	}
}

func handleGetRegister(w http.ResponseWriter, req *http.Request) {
	err := templates.Render(w, "register.gohtml", templates.Data{
		HeaderData: templates.HeaderData{
			Hide: true,
		},
		HideSidebar: true,
	})
	if err != nil {
		log.Println(err)
	}
}

func handlePostRegister(w http.ResponseWriter, req *http.Request) {
	data := templates.Data{
		HeaderData: templates.HeaderData{
			Hide:              true,
			IsUnauthenticated: true,
		},
		HideSidebar: true,
	}

	username := req.FormValue("username")
	if username == "" || strings.TrimSpace(username) == "" {
		data.FormErrorData.Username = "Username cannot be empty"
	}

	email := req.FormValue("email")
	if email == "" || !regex.Email.MatchString(email) {
		data.FormErrorData.Email = "Email is incorrect or already taken."
	}

	password := req.FormValue("password")
	if password == "" {
		data.FormErrorData.Password = "Password cannot be empty."
	}

	confirm := req.FormValue("confirm")
	if password != confirm {
		data.FormErrorData.Password = "Passwords do not match."
	}

	if !data.FormErrorData.IsEmpty() {
		w.WriteHeader(http.StatusBadRequest)
		err := templates.Render(w, "register.gohtml", data)
		if err != nil {
			log.Println(err)
		}
		return
	}

	u, err := config.App().Repo.CreateUser(username, email, password)
	if err != nil {
		data.FormErrorData.Username = err.Error()
		err := templates.Render(w, "register.gohtml", data)
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = setCookie(w, req, u, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// SignIn handles the [GET,POST] /auth/signin endpoint.
func SignIn(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		handleGetSignIn(w, req)
	case http.MethodPost:
		handlePostSignIn(w, req)
	}
}

func handleGetSignIn(w http.ResponseWriter, req *http.Request) {
	err := templates.Render(w, "signin.gohtml", templates.Data{
		HeaderData: templates.HeaderData{
			Hide:              true,
			IsUnauthenticated: true,
		},
		HideSidebar: true,
	})
	if err != nil {
		log.Println(err)
	}
}

func handlePostSignIn(w http.ResponseWriter, req *http.Request) {
	var errors templates.FormErrorData

	id := req.FormValue("username-or-email")
	if id == "" {
		errors.Username = "Username/email cannot be empty"
	}

	password := req.FormValue("password")
	if password == "" {
		errors.Password = "Password cannot be empty."
	}

	u := config.App().Repo.User(id)
	if !u.IsPasswordOk(password) {
		errors.Password = "Credentials are incorrect."
	}

	if !errors.IsEmpty() {
		w.WriteHeader(http.StatusBadRequest)
		err := templates.Render(w, "signin.gohtml", templates.Data{
			HeaderData: templates.HeaderData{
				Hide:              true,
				IsUnauthenticated: true,
			},
			HideSidebar:   true,
			FormErrorData: errors,
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	rememberMe := req.FormValue("remember-me")
	err := setCookie(w, req, u, rememberMe == "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func setCookie(w http.ResponseWriter, r *http.Request, u models.User, isSession bool) error {
	sid := uuid.NewString()
	token, err := auth.CreateToken(sid)
	if err != nil {
		log.Println("could not create token in login", err)
		return errors.New("our server didn't get enough lunch and is not working 200% right now. Try back later")
	}

	repository.Sessions[sid] = models.Session{
		UserID:       u.ID,
		UserInitials: u.GetInitials(),
	}

	c := http.Cookie{Name: "session", Value: token, Path: "/"}
	if !isSession {
		c.Expires = auth.ExpiresAt
	}
	http.SetCookie(w, &c)

	return nil
}

// SignOut handles the POST /auth/signout endpoint.
func SignOut(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		c = &http.Cookie{Name: "session", Value: ""}
	}
	c.MaxAge = -1

	sid, err := auth.ParseToken(c.Value)
	if err != nil {
		http.Error(w, "error parsing token from cookie", http.StatusInternalServerError)
		return
	}
	delete(repository.Sessions, sid)

	http.SetCookie(w, c)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
