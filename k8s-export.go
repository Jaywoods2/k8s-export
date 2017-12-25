package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ghodss/yaml"

	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const ROOT_DIR = "export/"
const SYSTEM_ROLE = "system:*|kubeadm:*"

func main() {
	var clientset = buildClient()
	//pv
	getPvs(clientset)

	var ns, _ = getNameSpaces(clientset)
	for _, v := range ns {
		namespace := v.GetName()
		if namespace == "default" ||
			namespace == "kube-system" ||
			namespace == "kube-public" {
			continue
		}
		fmt.Println(namespace)
		//pvc
		getPvcs(clientset, namespace)
		//deploy
		getDeploy(clientset, namespace)
		//svc
		getSvc(clientset, namespace)
		//ingress
		getIngress(clientset, namespace)
		//statefulset
		getStateFulSet(clientset, namespace)
		//sa
		getSas(clientset, namespace)
		//role
		getRoles(clientset, namespace)
		//rolebinding
		getRoleBinds(clientset, namespace)
		//cm
		getCms(clientset, namespace)
		//secret
		getSecrets(clientset, namespace)
	}

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func buildClient() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func getPvs(client *kubernetes.Clientset) {

	dir := ROOT_DIR + "pv/"
	fmt.Println("      |----pv")
	dataList, _ := client.CoreV1().PersistentVolumes().List(meta_v1.ListOptions{})
	for _, v := range dataList.Items {
		v.Kind = "PersistentVolume"
		v.APIVersion = "v1"
		v.CreationTimestamp.Reset()
		v.ResourceVersion = ""
		v.SelfLink = ""
		v.SetUID("")
		v.Annotations = nil
		v.Spec.ClaimRef.Reset()
		v.Status.Reset()
		fmt.Println("      |      |----" + v.GetName())
		os.MkdirAll(dir, os.ModePerm)
		data, _ := yaml.Marshal(v)
		err := ioutil.WriteFile(dir+v.GetName()+".yaml", data, os.ModePerm)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

func getPvcs(client *kubernetes.Clientset, namespace string) {
	dir := ROOT_DIR + "namespace/" + namespace + "/pvc/"
	dataList, _ := client.CoreV1().PersistentVolumeClaims(namespace).List(meta_v1.ListOptions{})
	if len(dataList.Items) == 0 {
		return
	}
	fmt.Println("      |----pvc")
	for _, v := range dataList.Items {
		fmt.Println("      |      |----" + v.GetName())
		v.Kind = "PersistentVolumeClaim"
		v.APIVersion = "v1"
		v.CreationTimestamp.Reset()
		v.ResourceVersion = ""
		v.SelfLink = ""
		v.SetUID("")
		v.Annotations = nil
		v.Spec.VolumeName = ""
		v.Status.Reset()

		os.MkdirAll(dir, os.ModePerm)
		data, _ := yaml.Marshal(v)
		err := ioutil.WriteFile(dir+v.GetName()+".yaml", data, os.ModePerm)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

func getNameSpaces(client *kubernetes.Clientset) ([]v1.Namespace, error) {
	ns, err := client.CoreV1().Namespaces().List(meta_v1.ListOptions{})
	for _, namespace := range ns.Items {
		if namespace.GetName() == "default" ||
			namespace.GetName() == "kube-system" ||
			namespace.GetName() == "kube-public" {
			continue
		}
		namespace.APIVersion = "v1"
		namespace.Kind = "Namespace"
		namespace.Spec.Reset()
		namespace.Status.Reset()
		namespace.ObjectMeta.SetUID("")
		namespace.ObjectMeta.SetSelfLink("")
		namespace.ObjectMeta.SetResourceVersion("")
		namespace.CreationTimestamp.Reset()

		yamlData, err := yaml.Marshal(namespace)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		directory := ROOT_DIR + "namespace/" + namespace.GetName()
		os.MkdirAll(directory, os.ModePerm)
		err = ioutil.WriteFile(directory+"/namespace.yaml", yamlData, os.ModePerm)

		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
	return ns.Items, err
}

func getDeploy(client *kubernetes.Clientset, namespace string) {
	deployDirectory := ROOT_DIR + "namespace/" + namespace + "/deploy/"

	deployList, _ := client.AppsV1beta1().Deployments(namespace).List(meta_v1.ListOptions{})
	if len(deployList.Items) == 0 {
		return
	}
	fmt.Println("      |----deploy")
	for _, deploy := range deployList.Items {
		fmt.Println("      |      |----" + deploy.GetName())
		deploy.Kind = "Deployment"
		deploy.APIVersion = "extensions/v1beta1"
		deploy.CreationTimestamp.Reset()
		deploy.ResourceVersion = ""
		deploy.SelfLink = ""
		deploy.SetUID("")
		deploy.Annotations = nil
		deploy.Status.Reset()
		os.MkdirAll(deployDirectory, os.ModePerm)
		deployData, _ := yaml.Marshal(deploy)
		err := ioutil.WriteFile(deployDirectory+deploy.GetName()+".yaml", deployData, os.ModePerm)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

func getSvc(client *kubernetes.Clientset, namespace string) {
	serviceList, err := client.CoreV1().Services(namespace).List(meta_v1.ListOptions{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	svcDirectory := ROOT_DIR + "namespace/" + namespace + "/svc/"
	if len(serviceList.Items) == 0 {
		return
	}
	fmt.Println("      |----svc")
	for _, svc := range serviceList.Items {
		fmt.Println("      |      |----" + svc.GetName())
		svc.Kind = "Service"
		svc.APIVersion = "v1"
		svc.CreationTimestamp.Reset()
		svc.ResourceVersion = ""
		svc.SelfLink = ""
		svc.SetUID("")
		svc.Annotations = nil
		svc.Spec.ClusterIP = ""
		svc.Status.Reset()
		os.MkdirAll(svcDirectory, os.ModePerm)
		svcData, _ := yaml.Marshal(svc)
		err = ioutil.WriteFile(svcDirectory+svc.GetName()+".yaml", svcData, os.ModePerm)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

func getIngress(client *kubernetes.Clientset, namespace string) {
	ingressList, _ := client.ExtensionsV1beta1().Ingresses(namespace).List(meta_v1.ListOptions{})
	dir := ROOT_DIR + "namespace/" + namespace + "/ingress/"
	if len(ingressList.Items) == 0 {
		return
	}
	fmt.Println("      |----ingress")
	for _, v := range ingressList.Items {
		fmt.Println("      |      |----" + v.GetName())
		v.Kind = "Ingress"
		v.APIVersion = "extensions/v1beta1"
		v.CreationTimestamp.Reset()
		v.ResourceVersion = ""
		v.SelfLink = ""
		v.SetUID("")
		v.Annotations = nil
		v.Status.Reset()
		os.MkdirAll(dir, os.ModePerm)
		data, _ := yaml.Marshal(v)
		err := ioutil.WriteFile(dir+v.GetName()+".yaml", data, os.ModePerm)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

func getStateFulSet(client *kubernetes.Clientset, namespace string) {
	dir := ROOT_DIR + "namespace/" + namespace + "/statefulset/"

	dataList, _ := client.AppsV1beta1().StatefulSets(namespace).List(meta_v1.ListOptions{})
	if len(dataList.Items) == 0 {
		return
	}
	fmt.Println("      |----statefulset")
	for _, v := range dataList.Items {
		fmt.Println("      |      |----" + v.GetName())
		v.Kind = "StatefulSet"
		v.APIVersion = "apps/v1beta1"
		v.CreationTimestamp.Reset()
		v.ResourceVersion = ""
		v.SelfLink = ""
		v.SetUID("")
		v.Annotations = nil
		v.Status.Reset()

		os.MkdirAll(dir, os.ModePerm)
		data, _ := yaml.Marshal(v)
		err := ioutil.WriteFile(dir+v.GetName()+".yaml", data, os.ModePerm)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

func getSas(client *kubernetes.Clientset, namespace string) {
	dir := ROOT_DIR + "namespace/" + namespace + "/sa/"

	dataList, err := client.CoreV1().ServiceAccounts(namespace).List(meta_v1.ListOptions{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(dataList.Items) == 1 && dataList.Items[0].GetName() == "default" {
		return
	}
	fmt.Println("      |----sa")
	for _, v := range dataList.Items {
		if v.GetName() == "default" {
			continue
		}
		fmt.Println("      |      |----" + v.GetName())
		v.Kind = "ServiceAccount"
		v.APIVersion = "v1"
		v.CreationTimestamp.Reset()
		v.ResourceVersion = ""
		v.SelfLink = ""
		v.SetUID("")
		v.Secrets = []v1.ObjectReference{}
		v.Annotations = nil
		os.MkdirAll(dir, os.ModePerm)
		data, _ := yaml.Marshal(v)
		err := ioutil.WriteFile(dir+v.GetName()+".yaml", data, os.ModePerm)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

func getRoles(client *kubernetes.Clientset, namespace string) {
	dir := ROOT_DIR + "namespace/" + namespace + "/role/"

	dataList, err := client.RbacV1().Roles(namespace).List(meta_v1.ListOptions{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(dataList.Items) == 0 {
		return
	}
	fmt.Println("      |----role")
	for _, v := range dataList.Items {

		match, _ := regexp.Match(SYSTEM_ROLE, []byte(v.GetName()))
		if match {
			continue
		}
		fmt.Println("      |      |----" + v.GetName())
		v.Kind = "Role"
		v.APIVersion = "rbac.authorization.k8s.io/v1"
		v.Annotations = nil
		v.CreationTimestamp.Reset()
		v.ResourceVersion = ""
		v.SelfLink = ""
		v.SetUID("")

		os.MkdirAll(dir, os.ModePerm)
		data, _ := yaml.Marshal(v)
		err := ioutil.WriteFile(dir+v.GetName()+".yaml", data, os.ModePerm)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

func getRoleBinds(client *kubernetes.Clientset, namespace string) {
	dir := ROOT_DIR + "namespace/" + namespace + "/rolebinding/"
	dataList, err := client.RbacV1().Roles(namespace).List(meta_v1.ListOptions{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(dataList.Items) == 0 {
		return
	}
	fmt.Println("      |----rolebinding")
	for _, v := range dataList.Items {
		fmt.Println("      |      |----" + v.GetName())
		v.Kind = "RoleBinding"
		v.APIVersion = "rbac.authorization.k8s.io/v1"
		v.Annotations = nil
		v.CreationTimestamp.Reset()
		v.ResourceVersion = ""
		v.SelfLink = ""
		v.SetUID("")
		os.MkdirAll(dir, os.ModePerm)
		data, _ := yaml.Marshal(v)
		err := ioutil.WriteFile(dir+v.GetName()+".yaml", data, os.ModePerm)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

func getCms(client *kubernetes.Clientset, namespace string) {
	dir := ROOT_DIR + "namespace/" + namespace + "/cm/"
	dataList, err := client.CoreV1().ConfigMaps(namespace).List(meta_v1.ListOptions{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(dataList.Items) == 0 {
		return
	}
	fmt.Println("      |----cm")
	for _, v := range dataList.Items {
		fmt.Println("      |      |----" + v.GetName())
		v.Kind = "ConfigMap"
		v.APIVersion = "v1"
		v.Annotations = nil
		v.CreationTimestamp.Reset()
		v.ResourceVersion = ""
		v.SelfLink = ""
		v.SetUID("")
		os.MkdirAll(dir, os.ModePerm)
		data, _ := yaml.Marshal(v)
		err := ioutil.WriteFile(dir+v.GetName()+".yaml", data, os.ModePerm)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

func getSecrets(client *kubernetes.Clientset, namespace string) {
	dir := ROOT_DIR + "namespace/" + namespace + "/secret/"
	dataList, err := client.CoreV1().Secrets(namespace).List(meta_v1.ListOptions{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(dataList.Items) == 0 {
		return
	}
	fmt.Println("      |----secret")
	for _, v := range dataList.Items {
		if v.Type == "kubernetes.io/service-account-token" {
			continue
		}
		fmt.Println("      |      |----" + v.GetName())
		v.Kind = "Secret"
		v.APIVersion = "v1"
		v.Annotations = nil
		v.CreationTimestamp.Reset()
		v.ResourceVersion = ""
		v.SelfLink = ""
		v.SetUID("")
		os.MkdirAll(dir, os.ModePerm)
		data, _ := yaml.Marshal(v)
		err := ioutil.WriteFile(dir+v.GetName()+".yaml", data, os.ModePerm)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}
