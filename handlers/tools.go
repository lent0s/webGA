package handlers

import (
	"fmt"
	"net/http"
)

func NotSupportedMethod(request, support string, w http.ResponseWriter) bool {

	if request != support {
		w.Header().Set("Support", support)
		err := fmt.Sprintf("method \"%s\" not allowed (%s)", request, support)
		http.Error(w, err, http.StatusMethodNotAllowed)
		return true
	}
	return false
}
