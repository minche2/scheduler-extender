package main
import (

	"k8s.io/api/core/v1"
    extenderv1 "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

type Prioritize struct {
	Name string
	Func func(pod v1.Pod, nodes []v1.Node) (*extenderv1.HostPriorityList, error)
}

func(prioitize Prioritize) Handler(args extenderv1.ExtenderArgs) *extenderv1.HostPriorityList {
	hostPriorityList, err := prioitize.Func(*args.Pod, args.Nodes.Items)
	if err != nil {
		return nil
	} else {
		return hostPriorityList
	}
}