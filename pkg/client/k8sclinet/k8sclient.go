package k8sclinet

import (
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/golang/glog"
	"flag"
)

/**
GetKubetnetesClient is get kubernetes client
 */
func GetKubetnetesClient() (kubeClient *kubernetes.Clientset, err error) {

	kubeconfig := flag.String("kubeconfig", "config", "absolute path to the kubeconfig file")

	config, errConfig := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if errConfig != nil {

		glog.Errorf("get kubetnetes config failed %#v\n",errConfig)
		return kubeClient, errConfig
	}

	kubeClient, err = kubernetes.NewForConfig(config)

	if err != nil {
		glog.Errorf("get kubernetes client failed %#v\n",err)
		return kubeClient, err
	}

	return kubeClient, nil

}