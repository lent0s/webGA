package cmd

import (
	"github.com/lent0s/webGA/handlers"
	"net/http"
)

func (p *Page) home(w http.ResponseWriter, r *http.Request) {

	if handlers.NotSupportedMethod(r.Method, http.MethodGet, w) {
		return
	}

	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	runTemplate(p, w)
}

func getSettings(w http.ResponseWriter, r *http.Request) {

	if handlers.NotSupportedMethod(r.Method, http.MethodPost, w) {
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := getAccountInfo(&r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func getStateInstance(w http.ResponseWriter, r *http.Request) {

	if handlers.NotSupportedMethod(r.Method, http.MethodPost, w) {
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := getAccountState(&r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {

	if handlers.NotSupportedMethod(r.Method, http.MethodPost, w) {
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := sendMessageT(&r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func sendFileByUrl(w http.ResponseWriter, r *http.Request) {

	if handlers.NotSupportedMethod(r.Method, http.MethodPost, w) {
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := sendFileByUrlT(&r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
