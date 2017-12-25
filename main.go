package main

import (
	"fmt"

	"github.com/Jaywoods2/k8s-export/config"
	"k8s.io/kubernetes/pkg/kubelet/kubeletconfig/util/log"
)

var cluster = config.ClusterConfiguration{
	Clientset: config.ClientSet,
}

func main() {
	ns, err := cluster.Namespaces()
	if err != nil {
		panic(err.Error())
	}
	for _, v := range ns {
		namespace := v.GetName()
		if namespace == "default" || namespace == "kube-system" || namespace == "kube-public" {
			continue
		}
		deploys, err := cluster.Deploys(namespace)
		if err != nil {
			log.Errorf("%s", err)
		}
		for _, v := range deploys {

		}
		fmt.Println(namespace)
	}
}
