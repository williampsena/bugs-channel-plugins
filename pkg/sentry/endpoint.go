// This package includes sentry plugin integration
package sentry

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/williampsena/bugs-channel-plugins/pkg/event"
)

type EndpointHandler func(w http.ResponseWriter, req *http.Request)

func HealthCheckEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Keep calm I'm absolutely alive 🐛")
}

func PostEventEndpoint(dispatcher event.EventsDispatcher) EndpointHandler {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		events, err := TranslateEvents(vars["id"], req.Body)

		if err != nil {
			HandleErrors(w, err, http.StatusUnprocessableEntity)
			return
		}

		dispatcher.DispatchMany(events)

		w.WriteHeader(http.StatusNoContent)
	}
}

func NoRouteEndpoint(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Oops! 👀")
}

func HandleErrors(w http.ResponseWriter, err error, httpStatus int) {
	w.WriteHeader(httpStatus)
	log.Errorf("⛔ %v", err.Error())
	fmt.Fprintf(w, "%v", err.Error())
}
