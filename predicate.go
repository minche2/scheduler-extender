package main

import (
	"k8s.io/api/core/v1"
	extenderv1 "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)


type Predicate struct {
	Name string
	Func func(pod v1.Pod, node v1.Node) (bool, error)
}

func(p Predicate) Handler(args extenderv1.ExtenderArgs) *extenderv1.ExtenderFilterResult {
	pod := args.Pod
	canSchedule := make([]v1.Node, 0, len(args.Nodes.Items))
	cannotSchedule := make(map[string]string)
	//canScheduleName := make([]string, len(args.Nodes.Items))

	for _, node := range args.Nodes.Items {
		result, err := p.Func(*pod, node)
		if err != nil {
			cannotSchedule[node.Name] = err.Error()
		} else{
			if result {
				canSchedule = append(canSchedule, node)
				//canScheduleName = append(canScheduleName, node.Name)
			}
		}
	}

	result := extenderv1.ExtenderFilterResult {
		Nodes: &v1.NodeList{
			Items: canSchedule,
		},
		FailedNodes: cannotSchedule,
		Error: "",
	}

	return &result
}