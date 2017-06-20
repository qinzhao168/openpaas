package app

import (
	"openpaas/pkg/client/k8sclinet"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	v1beta11 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"net/http"
	"fmt"
	//"k8s.io/kubernetes/pkg/kubectl"
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
	serport := v1.ServicePort{
		Name:"mysqlport",
		Protocol:v1.ProtocolTCP,
		Port:7878,
	}
	service, err := kc.KubeCtl.CoreV1().Services("qinzhao").Create(&v1.Service{
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
	if err != nil {
		glog.Errorf("%#v\n", err)
		return
	}
	glog.Info(service)
	return

}

func (kc *KubetnetesCLient)CreateDeployment() {

	Container := v1.Container{
		Name:"redis",
		Image:"redis:latest",
	}

	dep := &v1beta1.Deployment{

		TypeMeta:metav1.TypeMeta{
			Kind:"Deployment",
			APIVersion:"v1beta1",
		},
		ObjectMeta:metav1.ObjectMeta{
			Name:"qinzhao",
			Namespace:"qinzhao",
		},
		Spec:v1beta1.DeploymentSpec{
			Replicas:Int32Toint32Point(1),
			Selector:&metav1.LabelSelector{
				MatchLabels:map[string]string{
					"qinzhao":"qinzhao",
				},
			},
			Template:v1.PodTemplateSpec{
				ObjectMeta:metav1.ObjectMeta{
					Name:"qinzhao",
					Namespace:"qinzhao",
					Labels:map[string]string{
						"qinzhao":"qinzhao",
					},
				},
				Spec:v1.PodSpec{
					Containers:[]v1.Container{Container},
				},
			},

		},
	}

	deployment1, err := kc.KubeCtl.AppsV1beta1Client.Deployments("qinzhao").Create(dep)
	if err != nil {
		glog.Errorf("create deployment failed%#v\n", err)
		return
	}

	fmt.Println(deployment1)
	return
}

func (kc *KubetnetesCLient)CreateDaemonSets() {
	Container := v1.Container{
		Name:"redis",
		Image:"redis:latest",
		ImagePullPolicy:v1.PullAlways,
	}
	daemonSets := &v1beta11.DaemonSet{
		TypeMeta:metav1.TypeMeta{
			Kind:"DaemonSet",
			APIVersion:"v1beta1",
		},
		ObjectMeta:metav1.ObjectMeta{
			Name:"qinzhao",
			Namespace:"qinzhao",
			Labels:map[string]string{
				"qinzhao":"qinzhao",
			},
		},
		Spec:v1beta11.DaemonSetSpec{
			Selector:&metav1.LabelSelector{
				MatchLabels:map[string]string{"qinzhao":"qinzhao"},
			},
			Template:v1.PodTemplateSpec{
				ObjectMeta:metav1.ObjectMeta{
					Namespace:"qinzhao",
					Name:"qinzhao",
					Labels:map[string]string{
						"qinzhao":"qinzhao",
					},
				},
				Spec:v1.PodSpec{
					Containers:[]v1.Container{Container},
				},
			},
		},

	}

	daemonset, err := kc.KubeCtl.DaemonSets("qinzhao").Create(daemonSets)
	if err != nil {
		glog.Errorf("%#s\n", err)
		return
	}
	fmt.Println(daemonset)
	return

}

func (kc *KubetnetesCLient) GetDaemonSets() {
	//opention:=v1.
	daemonset, err := kc.KubeCtl.DaemonSets("qinzhao").Get("qinzhao", metav1.GetOptions{
		TypeMeta:metav1.TypeMeta{
			Kind:"DaemonSet",
			APIVersion:"v1beta1",
		},
	})
	if err != nil {
		glog.Errorf("get daemonSets failed", err)
		return
	}
	fmt.Println(daemonset.Name)
	return

}

func (kc *KubetnetesCLient)GetService() {
	service, err := kc.KubeCtl.CoreV1().Services("qinzhao").Get("qinzhao", metav1.GetOptions{
		TypeMeta:metav1.TypeMeta{
			Kind:"Service",
			APIVersion:"v1",
		},
	})
	if err != nil {
		glog.Errorf("get service failed%#v\n", service)
		return
	}
	fmt.Println(service.Spec.ClusterIP)
}

func (kc *KubetnetesCLient)GetDeployment() {
	deployment, err := kc.KubeCtl.AppsV1beta1Client.Deployments("qinzhao").Get("qinzhao", metav1.GetOptions{
		TypeMeta:metav1.TypeMeta{
			Kind:"Deployment",
			APIVersion:"v1beta1",
		},
	})

	if err != nil {
		glog.Errorf("get deploymnet failed%#v\n", err)

		return
	}

	fmt.Println(deployment.Spec.Template)
	return
}

func (kc *KubetnetesCLient)GetRC() {

	rc, err := kc.KubeCtl.CoreV1().ReplicationControllers("qinzhao").Get("qinzhao",
		metav1.GetOptions{
			TypeMeta:metav1.TypeMeta{
				Kind:"ReplicationController",
				APIVersion:"v1",
			},
		})
	if err != nil {
		glog.Errorf("%#v\n", err)
		return
	}

	fmt.Println(rc.ResourceVersion)

}

func (kc *KubetnetesCLient)UpdateNamespace() {
	namespace := &v1.Namespace{
		TypeMeta:metav1.TypeMeta{
			Kind:"Namespace",
			APIVersion:"v1",
		},
		ObjectMeta:metav1.ObjectMeta{
			Name:"qinzhao",
			Namespace:"qinzhao",
			Annotations:map[string]string{
				"qinzhao":"qinzhao",
			},
		},


	}
	name, err := kc.KubeCtl.CoreV1().Namespaces().Update(namespace)
	if err != nil {
		glog.Errorf("update namespace failed %#v\n", err)
		return
	}
	fmt.Println(name.ResourceVersion)
	return
}

func (kc *KubetnetesCLient)GetNameSpace() {

	namespace, err := kc.KubeCtl.CoreV1().Namespaces().Get("qinzhao", metav1.GetOptions{
		TypeMeta:metav1.TypeMeta{
			Kind:"Namespace",
			APIVersion:"v1",
		},

	})
	if err != nil {
		glog.Errorf("get namespace failed%#v\n", err)
		return
	}
	fmt.Println(namespace.ResourceVersion)
	return

}

func (kc *KubetnetesCLient)UpdateRC(){
	Container := v1.Container{
		Name:"redis",
		Image:"redis:v1.2",
		ImagePullPolicy:v1.PullAlways,
	}
	rc,err:=kc.KubeCtl.CoreV1().ReplicationControllers("qinzhao").Update(&v1.ReplicationController{
		TypeMeta:metav1.TypeMeta{
			Kind:"ReplicationController",
			APIVersion:"v1",
		},
		ObjectMeta:metav1.ObjectMeta{
			Namespace:"qinzhao",
			Name:"qinzhao",
			Annotations:map[string]string{
			"qinzhao":"qinzhao",
			},

		},
		Spec:v1.ReplicationControllerSpec{
			Selector:map[string]string{
				"qinzhao":"qinzhao",
			},
			Template:&v1.PodTemplateSpec{
				ObjectMeta:metav1.ObjectMeta{
					Name:"qinzhao",
					Namespace:"qinzhao",
					Labels:map[string]string{
						"qinzhao":"qinzhao",
					},
				},
				Spec:v1.PodSpec{
					Containers:[]v1.Container{Container},
				},
			},
		},
	})

	if err!=nil{
		glog.Errorf("update ReplicationController failed %#v\n",err)
		return
	}
	fmt.Println(rc.ResourceVersion)
	return
}

//func (kc *KubetnetesCLient)UpdateDaemonSets(){
//kc.KubeCtl
//
//}



func Run() {

	var kubecli KubetnetesCLient

	kubecli.getKubeCli()
	//kubecli.GetNodes()
	//kubecli.CreateNamespace()
	//kubecli.GetPods()
	//kubecli.CreateRC()
	//kubecli.CreateDeployment()
	//kubecli.CreateService()
	//kubecli.CreateDaemonSets()

	//kubecli.GetDaemonSets()

	//kubecli.GetService()
	//kubecli.GetDeployment()

	//kubecli.GetRC()
	//kubecli.GetNameSpace()
	//kubecli.UpdateNamespace()
	//kubecli.GetNameSpace()
	kubecli.GetRC()
	kubecli.UpdateRC()
	http.ListenAndServe(fmt.Sprintf(":%d", 9090), nil)

}


