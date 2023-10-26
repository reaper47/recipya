package server

import (
	"fmt"
	"net/http"
)

func (s *Server) integrationsImportNextcloud(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	baseURL := r.FormValue("url")
	if username == "" || password == "" || baseURL == "" {
		w.Header().Set("HX-Trigger", makeToast("Invalid username, password or URL.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipes, err := s.Integrations.NextcloudImport(http.DefaultClient, baseURL, username, password, s.Files)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Failed to import Nextcloud recipes.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := getUserID(r)

	count := 0
	skipped := 0
	for _, r := range *recipes {
		_, err := s.Repository.AddRecipe(&r, userID)
		if err != nil {
			skipped++
			continue
		}
		count++
	}

	w.Header().Set("HX-Trigger", makeToast(fmt.Sprintf("Imported %d recipes. Skipped %d.", count, skipped), infoToast))
	w.WriteHeader(http.StatusCreated)
}
