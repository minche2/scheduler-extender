package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	extenderv1 "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
	"log"
	"net/http"
)

func DebugLogging(h httprouter.Handle, path string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Print("debug: ", path, " request body = ", r.Body)
		h(w, r, p)
		log.Print("debug: ", path, " response=", w)
	}
}

func checkBody(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
}

func versionRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, fmt.Sprint(version))
}

func predicateRoute(p Predicate) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		checkBody(w, r)

		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)
		log.Print("info: ", p.Name, " ExtenderArgs = ", buf.String())

		var extenderArgs extenderv1.ExtenderArgs
		var extenderFilterResult *extenderv1.ExtenderFilterResult

		if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
			extenderFilterResult = &extenderv1.ExtenderFilterResult{
				Nodes: nil,
				FailedNodes: nil,
				Error: err.Error(),
				NodeNames: nil,
			}
		} else {
			extenderFilterResult = p.Handler(extenderArgs)
		}

		if resultBody, err := json.Marshal(extenderFilterResult); err != nil {
			panic(err)
		} else {
			log.Print("info: ", p.Name, " extenderFilterResult = ", string(resultBody))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resultBody)
		}
	}
}

func prioritizeRoute(prioritize Prioritize) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		checkBody(w, r)

		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)
		log.Print("info: ", prioritize.Name, " ExtenderArgs = ", buf.String())

		var extenderArgs extenderv1.ExtenderArgs
		var HostPriorityList *extenderv1.HostPriorityList

		if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
			HostPriorityList = &extenderv1.HostPriorityList{}
		} else{
			HostPriorityList = prioritize.Handler(extenderArgs)
		}

		if resultBody, err := json.Marshal(HostPriorityList); err != nil {
			panic(err)
		} else {
			log.Print("info: ", prioritize.Name, " HostPriorityList = ", string(resultBody))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resultBody)
		}

	}
}

func AddVersion(router *httprouter.Router) {
	router.GET(versionPath, DebugLogging(versionRoute, versionPath))
}

func AddPredicate(router *httprouter.Router, predicate Predicate){
	path := predicatesPrefix + "/" + predicate.Name
	router.POST(path, DebugLogging(predicateRoute(predicate), path))
}

func AddPriority(router *httprouter.Router, priority Prioritize) {
	path := prioritiesPrefix + "/" + priority.Name
	router.POST(path, DebugLogging(prioritizeRoute(priority), path))
}