package main

import (
	"github.com/comail/colog"
	"github.com/julienschmidt/httprouter"
	"k8s.io/api/core/v1"
	extenderv1 "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	version = "0.0.1"
	versionPath = "/version"
	apiPrefix = "/scheduler"
	//bindPath = apiPrefix + "/bind"
	predicatesPrefix = apiPrefix + "/predicates"
	prioritiesPrefix = apiPrefix + "/priorities"
)

var (
	firstPredicate = Predicate{
		Name: "all-true",
		Func: func(pod v1.Pod, node v1.Node) (b bool, err error) {
			return true, nil
		},
	}

	secondPredicate = Predicate{
		Name: "second-all-true",
		Func: func(pod v1.Pod, node v1.Node) (b bool, err error) {
			return true, nil
		},
	}

	//extender里返回 *extenderv1.HostPriorityList，包含了nodeName与Score的一一对应
	firstPriority = Prioritize{
		Name: "zero-score",
		Func: func(_ v1.Pod, nodes []v1.Node) (*extenderv1.HostPriorityList, error) {
			var priorityList extenderv1.HostPriorityList
			priorityList = make([]extenderv1.HostPriority, len(nodes))
			for i, node := range nodes {
				priorityList[i] = extenderv1.HostPriority {
					Host: node.Name,
					Score: 0,
				}
			}
			return &priorityList, nil
		},
	}
)

func stringToLevel(levelStr string) colog.Level {
	level := strings.ToUpper(levelStr)
	switch level {
	case "TRACE":
		return colog.LTrace
	case "DEBUG":
		return colog.LDebug
	case "INFO":
		return colog.LInfo
	case "WARNING":
		return colog.LWarning
	case "ERROR":
		return colog.LError
	case "ALERT":
		return colog.LAlert
	default:
		log.Printf("warning: LOG_LEVEL \"%s\" is not set or invalid, use default log level \"INFO\"", level)
		return colog.LInfo
	}
}

func main(){
	// set log level
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.SetFormatter(&colog.StdFormatter{
		Flag:        log.LstdFlags | log.Lshortfile,
		Colors:      true,
	})

	colog.Register()
	level := stringToLevel(os.Getenv("LOG_LEVEL"))
	log.Print("Log level is set to ", strings.ToUpper(level.String()))
	colog.SetMinLevel(level)

	//add router
	router := httprouter.New()   //生成了一个*Router路由指针, Router 实际上 is a http.Handler
	AddVersion(router)

	predicates := []Predicate{firstPredicate, secondPredicate}
	for _, predicate := range predicates{
		AddPredicate(router, predicate)
	}

	priorities := []Prioritize{firstPriority}
	for _, priority := range priorities{
		AddPriority(router, priority)
	}

	//start http server
	log.Print("info: server starting on the port :80")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)   // Fatal is equivalent to Print() followed by a call to os.Exit(1).
	}

}