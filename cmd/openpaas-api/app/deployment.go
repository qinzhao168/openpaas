package app

//import (
//	v1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
//	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/client-go/kubernetes"
//)
//
//
//type Deploymentor interface {
//	Create(namespaces string,deployment *v1beta1.Deployment)(*v1beta1.Deployment,error)
//	Update(namespaces string,deployment *v1beta1.Deployment)(*v1beta1.Deployment,error)
//	UpdateStatus(namespaces string,deployment *v1beta1.Deployment)(*v1beta1.Deployment,error)
//	Delete(name string,options v1.DeleteOptions)error
//	Get(name string,options v1.GetOptions)(*v1beta1.Deployment,error)
//	List(options v1.GetOptions)(*v1beta1.DeploymentList,error)
//}
//
//type deployment struct{
//	kubeclient	*kubernetes.Clientset
//	ns string
//}
//
//func newDeployments()*deployment{
//	return &deployment{
//		kubeclient:KubeClient,
//	}
//}

//func (c *deployment)Create(namespaces string,deployment *v1beta1.Deployment)(*v1beta1.Deployment,error){
//	c.kubeclient.Apps().Deployments(namespaces)
//}



