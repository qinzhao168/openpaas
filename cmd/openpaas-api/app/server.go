package app

import (
	"openpaas/pkg/client/k8sclinet"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"
	"net/http"
	"fmt"
)

type KubetnetesCLient struct {
	KubeCtl *kubernetes.Clientset
}
type TypeMeta struct {
	TypeMeta metav1.TypeMeta `json:",inline"`
}

type ObjectMeta struct {
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

func (kc *KubetnetesCLient)getKubeCli() {
	kubeclient, err := k8sclinet.GetKubetnetesClient()
	if err != nil {
		glog.Errorf("get k8s client failed%#v\n", err)
		return
	}
	if kubeclient != nil {
		kc.KubeCtl = kubeclient
		kubeclient = nil
	}
}

func (kc *KubetnetesCLient)GetNodes() (*v1.NodeList, error) {
	nodes, err := kc.KubeCtl.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("get all namespaces failed %#v\n", err)
		return nodes, err
	}
	glog.Info(nodes)
	return nodes, nil
}

func (kc *KubetnetesCLient)CreateNamespace() (*v1.Namespace, error) {

	namespace := &v1.Namespace{
		TypeMeta:metav1.TypeMeta{
			Kind:"Namespace",
			APIVersion:"v1",
		},
		ObjectMeta:metav1.ObjectMeta{
			Name:"qinzhao",
		},
	}
	return kc.KubeCtl.CoreV1().Namespaces().Create(namespace)
}

func (kc *KubetnetesCLient)GetPods() {
	podlist, err := kc.KubeCtl.CoreV1().Pods("default").List(metav1.ListOptions{
		TypeMeta:metav1.TypeMeta{
			Kind:"Pod",
			APIVersion:"v1",
		},

	})
	if err != nil {
		glog.Errorln("get podlist failed", err)
		return
	}
	glog.Info(podlist)
	return
}

func Int32Toint32Point(input int32) *int32 {
	tmp := new(int32)
	*tmp = int32(input)
	return tmp

}
/**
CreateServer is create rc
 */
func (kc *KubetnetesCLient)CreateRC() {
	//Selector := make(map[string]string, 1)
	//Selector["qinzhao"] = "qinzhao"
	//Containers := make([]v1.Container, 1)
	conatiner := v1.Container{
		Name:"qinzhao",
		Image:"redis:latest",

	}
	//Containers = append(Containers, conatiner)
	rc := &v1.ReplicationController{
		TypeMeta:metav1.TypeMeta{
			Kind:"ReplicationController",
			APIVersion:"v1",
		},
		ObjectMeta:metav1.ObjectMeta{
			Name:"qinzhao",
			Namespace:"qinzhao",
		},
		Spec:v1.ReplicationControllerSpec{
			Replicas:Int32Toint32Point(1),
			Selector:map[string]string{
				"name":"qinzhao",
			},
			Template:&v1.PodTemplateSpec{
				ObjectMeta:metav1.ObjectMeta{
					Name:"qinzhao",
					Namespace:"qinzhao",
					Labels:map[string]string{
						"name":"qinzhao",
					},
				},
				Spec:v1.PodSpec{
					Containers:[]v1.Container{conatiner},
				},
			},
		},


	}
	rc, err := kc.KubeCtl.CoreV1Client.ReplicationControllers("qinzhao").
		Create(rc)
	if err != nil {
		glog.Errorf("create rc failed %#s", err)
		return
	}
	fmt.Println(rc.Namespace)
	return
}
func (kc *KubetnetesCLient)CreateService() {
	serport:=v1.ServicePort{
		Name:"mysqlport",
		Protocol:v1.ProtocolTCP,
		Port:7878,
	}
	service,err:=kc.KubeCtl.CoreV1().Services("qinzhao").Create(&v1.Service{
		TypeMeta:metav1.TypeMeta{
			Kind:"Service",
			APIVersion:"v1",
		},
		ObjectMeta:metav1.ObjectMeta{
			Name:"qinzhao",
			Namespace:"qinzhao",
		},
		Spec:v1.ServiceSpec{
			Selector:map[string]string{
				"name":"qinzhao",
			},
			Ports:[]v1.ServicePort{serport},
		},
	})
	if err!=nil{
		glog.Errorf("%#v\n",err)
		return
	}
	glog.Info(service)
	return

}

func Run() {
	var kubecli KubetnetesCLient
	kubecli.getKubeCli()
	//kubecli.GetNodes()
	kubecli.CreateNamespace()
	//kubecli.GetPods()
	kubecli.CreateRC()
	kubecli.CreateService()
	http.ListenAndServe(fmt.Sprintf(":%d", 9090), nil)

}


