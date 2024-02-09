package server

import (
	"errors"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/internal/utils/regex"
	"github.com/reaper47/recipya/web/components"
	"log"
	"maps"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *Server) changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if app.Config.Server.IsAutologin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	currentPassword := r.FormValue("password-current")
	newPassword := r.FormValue("password-new")
	if currentPassword == newPassword {
		w.Header().Set("HX-Trigger", makeToast("New password is same as current.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	confirmPassword := r.FormValue("password-confirm")
	if confirmPassword != newPassword {
		w.Header().Set("HX-Trigger", makeToast("Passwords do not match.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := getUserID(r)

	if app.Config.Server.IsDemo && s.Repository.UserID("demo@demo.com") == userID {
		w.Header().Set("HX-Trigger", makeToast("Your Facebook password has been changed.", infoToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.Repository.IsUserPassword(userID, currentPassword) {
		w.Header().Set("HX-Trigger", makeToast("Current password is incorrect.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error encoding your password.", errorToast))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = s.Repository.UpdatePassword(userID, hashPassword)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Failed to update password.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", makeToast("Password updated.", infoToast))
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) confirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	query := r.URL.Query()
	if query == nil {
		log.Printf("confirmHandler.Query() returned nil")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	token := query.Get("token")
	if token == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userID, err := auth.ParseToken(token)
	if err != nil {
		log.Printf("[error] confirmHandler.ParseToken (token: %s): %q", token, err)
		w.WriteHeader(http.StatusBadRequest)

		_ = components.SimplePage(templates.ErrorTokenExpired.Title, templates.ErrorTokenExpired.Content).Render(r.Context(), w)
		return
	}

	err = s.Repository.Confirm(userID)
	if err != nil {
		log.Printf("[error] confirmHandler.Confirm (token: %s): %q", token, err)
		w.WriteHeader(http.StatusNotFound)

		const content = `An error occurred when you requested to confirm your account.
				The problem has been forwarded to our team automatically. We will look into it and come
                back to you. We apologise for this inconvenience.`
		_ = components.SimplePage("Confirm Error", content).Render(r.Context(), w)
		return
	}

	_ = components.SimplePage("Success", "Your account has been confirmed.").Render(r.Context(), w)
}

func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if app.Config.Server.IsAutologin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userID := getUserID(r)
	if app.Config.Server.IsDemo && s.Repository.UserID("demo@demo.com") == userID {
		w.Header().Set("HX-Trigger", makeToast("Your savings account has been deleted.", errorToast))
		w.WriteHeader(http.StatusTeapot)
		return
	}

	err := s.Repository.DeleteUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.logoutHandler(w, r)
}

func (s *Server) forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if getUserIDFromSessionCookie(r) != -1 || getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken) != -1 {
		w.Header().Set("HX-Redirect", "/settings")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	_ = components.ForgotPasswordPage().Render(r.Context(), w)
}

func (s *Server) forgotPasswordPostHandler(w http.ResponseWriter, r *http.Request) {
	if getUserIDFromSessionCookie(r) != -1 || getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken) != -1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	email := r.FormValue("email")
	if s.Repository.IsUserExist(email) {
		userID := s.Repository.UserID(email)
		token, err := auth.CreateToken(map[string]any{"userID": userID}, 1*time.Hour)
		if err != nil {
			log.Printf("[error] forgotPasswordPostHandler.CreateToken: %q", err)
			w.Header().Set("HX-Trigger", makeToast("Forgot password failed.", errorToast))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		username := "user"
		split := strings.Split(email, "@")
		if len(split) > 0 {
			username = split[0]
		}

		data := templates.EmailData{
			Token:    token,
			UserName: username,
			URL:      app.Config.Address(),
		}

		err = s.Email.Send(email, templates.EmailForgotPassword, data)
		if err != nil {
			log.Printf("[error] forgotPasswordPostHandler.SendEmail (data: %+v): %q", data, err)
			s.Email.Queue(email, templates.EmailForgotPassword, data)

			const content = "The email could not be sent because the SendGrid daily sent email quota has been reached. " +
				"The action has been logged. The next batch of emails will be sent tomorrow. " +
				"You can sponsor the author of this project or buy him a coffee for him to have enough money to purchase the paid SendGrid plan to increase the limit. " +
				"You will find the details here: https://github.com/reaper47/heavy-metal-notifier?tab=readme-ov-file#sponsors."

			_ = components.SimplePage("Email Quota Reached", content).Render(r.Context(), w)
			return
		}
	}

	const content = "An email with instructions on how to reset your password has been sent to you. Please check your inbox and follow the provided steps to regain access to your account."
	_ = components.SimplePage("Password Reset Requested", content).Render(r.Context(), w)
}

func forgotPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query == nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = components.SimplePage(templates.ErrorTokenExpired.Title, templates.ErrorTokenExpired.Content).Render(r.Context(), w)
		return
	}

	userID, err := auth.ParseToken(query.Get("token"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = components.SimplePage(templates.ErrorTokenExpired.Title, templates.ErrorTokenExpired.Content).Render(r.Context(), w)
		return
	}

	_ = components.ForgotPasswordResetPage(strconv.FormatInt(userID, 10)).Render(r.Context(), w)
}

func (s *Server) forgotPasswordResetPostHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.FormValue("user-id")
	password := r.FormValue("password")
	confirm := r.FormValue("password-confirm")
	if userIDStr == "" || password == "" || password != confirm {
		w.Header().Set("HX-Trigger", makeToast("Password is invalid.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Printf("[error] forgotPasswordResetPostHandler.ParseInt for (%d): %q", userID, err)
		w.Header().Set("HX-Trigger", makeToast("Password is invalid.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashPassword, err := auth.HashPassword(password)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error encoding your password.", errorToast))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = s.Repository.UpdatePassword(userID, hashPassword)
	if err != nil {
		log.Printf("[error] forgotPasswordResetPostHandler.UpdatePassword for %d: %q", userID, err)
		w.Header().Set("HX-Trigger", makeToast("Updating password failed.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/auth/login")
	w.Header().Set("HX-Trigger", makeToast("Password updated.", infoToast))
	w.WriteHeader(http.StatusSeeOther)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	_ = components.LoginPage(app.Config.Server.IsDemo, app.Config.Server.IsNoSignups).Render(r.Context(), w)
}

func (s *Server) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if !regex.Email.MatchString(email) || password == "" {
		w.Header().Set("HX-Trigger", makeToast("Credentials are invalid.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID := s.Repository.VerifyLogin(email, password)
	if userID == -1 {
		w.Header().Set("HX-Trigger", makeToast("Credentials are invalid.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sid := uuid.New()
	if SessionData == nil {
		SessionData = make(map[uuid.UUID]int64)
	}
	SessionData[sid] = userID
	http.SetCookie(w, NewSessionCookie(sid.String()))

	if r.FormValue("remember-me") == "yes" {
		selector, validator := auth.GenerateSelectorAndValidator()
		http.SetCookie(w, NewRememberMeCookie(selector, validator))
		err := s.Repository.AddAuthToken(selector, validator, userID)
		if err != nil {
			log.Printf("[error] loginPostHandler.AddAuthToken: %q", err)
		}
	}

	redirectURI := "/"
	c, err := r.Cookie(cookieNameRedirect)
	if c != nil && !errors.Is(err, http.ErrNoCookie) {
		redirectURI = c.Value
	}

	w.Header().Set("HX-Redirect", redirectURI)
	w.WriteHeader(http.StatusSeeOther)
}

func guideLoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	_ = components.RegisterPage().Render(r.Context(), w)
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if app.Config.Server.IsAutologin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userID := getUserIDFromSessionCookie(r)
	if userID == -1 {
		userID = getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken)
	}

	if userID == -1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sessionCookie, err := r.Cookie(cookieNameSession)
	if sessionCookie != nil && !errors.Is(err, http.ErrNoCookie) {
		c := NewSessionCookie(sessionCookie.Value)
		c.MaxAge = -1
		http.SetCookie(w, c)
	}
	maps.DeleteFunc(SessionData, func(_ uuid.UUID, id int64) bool { return id == userID })

	rememberMeCookie, err := r.Cookie(cookieNameRememberMe)
	if rememberMeCookie != nil && !errors.Is(err, http.ErrNoCookie) {
		err = s.Repository.DeleteAuthToken(userID)
		if err != nil {
			log.Printf("[error] logoutHandler.DeleteAuthToken: %q", err)
		}

		selector, validator, _ := strings.Cut(rememberMeCookie.Value, ":")
		c := NewRememberMeCookie(selector, validator)
		c.MaxAge = -1
		http.SetCookie(w, c)
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func (s *Server) registerPostHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if !regex.Email.MatchString(email) || password != r.FormValue("password-confirm") {
		w.Header().Set("HX-Trigger", makeToast("Registration failed. User might be registered or password invalid.", errorToast))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	hashPassword, err := auth.HashPassword(password)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error encoding your password.", errorToast))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userID, err := s.Repository.Register(email, hashPassword)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Registration failed. User might be registered or password invalid.", errorToast))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	token, err := auth.CreateToken(map[string]any{"userID": userID}, 14*24*time.Hour)
	if err != nil {
		log.Printf("[error] registerPostHandler.CreateToken: %q", err)
		w.Header().Set("HX-Trigger", makeToast("Registration failed. User might be registered or password invalid.", errorToast))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	username := "user"
	split := strings.Split(email, "@")
	if len(split) > 0 {
		username = split[0]
	}

	data := templates.EmailData{
		Token:    token,
		UserName: username,
		URL:      app.Config.Address(),
	}

	err = s.Email.Send(email, templates.EmailIntro, data)
	if err != nil {
		log.Printf("[error] registerPostHandler.SendEmail (data: %+v): %q", data, err)
		s.Email.Queue(email, templates.EmailIntro, data)
	}

	w.Header().Set("HX-Redirect", "/auth/login")
	w.WriteHeader(http.StatusSeeOther)
}
