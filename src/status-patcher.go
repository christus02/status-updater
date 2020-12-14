package main

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog"

	"k8s.io/client-go/rest"
	//"k8s.io/client-go/tools/clientcmd"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"os"
)

var SERVICE_NAME = os.Getenv("SERVICE_NAME")
var SERVICE_NAMESPACE = os.Getenv("SERVICE_NAMESPACE")
var EXTERNAL_ENDPOINT_TYPE_ANNOTATION = os.Getenv("EXTERNAL_ENDPOINT_TYPE_ANNOTATION")
var ENDPOINT_ANNOTATION = os.Getenv("ENDPOINT_ANNOTATION")
var ENDPOINT = ""

var clientset *kubernetes.Clientset

func InitWatcher() {

	config, err := rest.InClusterConfig() // Using in-cluster-config serviceAccount will be used
	//config, err := clientcmd.BuildConfigFromFlags("", "/home/raghulc/.kube/config")
	if err != nil {
		klog.Error(err.Error())
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		klog.Error(err.Error())
	}
	factory := informers.NewSharedInformerFactory(clientset, 0)
	serviceInformer := factory.Core().V1().Services().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
	})
	klog.Info("Starting Service Informer")
	go serviceInformer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, serviceInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper

}

func onAdd(obj interface{}) {
	svc := obj.(*corev1.Service)
	if svc.ObjectMeta.Name == SERVICE_NAME && svc.ObjectMeta.Namespace == SERVICE_NAMESPACE {
		value, ok := svc.ObjectMeta.Annotations[ENDPOINT_ANNOTATION]
		if ok {
			ENDPOINT = value
			klog.Infof("Service: %s, Namespace: %s, ENDPOINT: %s", SERVICE_NAME, SERVICE_NAMESPACE, ENDPOINT)
			patchData := []byte(`{"status":{"loadBalancer": {"ingress": [{"hostname": "abc.ingress.com"}]}}}`)
			res, err := clientset.CoreV1().Services(SERVICE_NAMESPACE).Patch(context.TODO(), SERVICE_NAME, types.StrategicMergePatchType, patchData, metav1.PatchOptions{FieldManager: "my-controller"})
			if err != nil {
				klog.Error(err)
			}
			klog.Info(res)

		} else {
			klog.Errorf("Service: %s, Namespace: %s, ENDPOINT: X -> NOT FOUND <- X", SERVICE_NAME, SERVICE_NAMESPACE)
		}

	}
}

func onUpdate(obj interface{}, newObj interface{}) {
	newSvc := newObj.(*corev1.Service)
	if newSvc.ObjectMeta.Name == SERVICE_NAME && newSvc.ObjectMeta.Namespace == SERVICE_NAMESPACE {
		value, ok := newSvc.ObjectMeta.Annotations[ENDPOINT_ANNOTATION]
		if ok {
			ENDPOINT = value
			klog.Infof("Service: %s, Namespace: %s, ENDPOINT: %s", SERVICE_NAME, SERVICE_NAMESPACE, ENDPOINT)
			patchData := []byte(`{"status":{"loadBalancer": {"ingress": [{"hostname": "abc.ingress.com"}]}}}`)
			res, err := clientset.CoreV1().Services(SERVICE_NAMESPACE).Patch(context.TODO(), SERVICE_NAME, types.StrategicMergePatchType, patchData, metav1.PatchOptions{FieldManager: "my-controller"})
			if err != nil {
				klog.Error(err)
			}
			klog.Info(res)

		} else {
			klog.Errorf("Service: %s, Namespace: %s, ENDPOINT: X -> NOT FOUND <- X", SERVICE_NAME, SERVICE_NAMESPACE)
		}

	}

}

func main() {
	InitWatcher()
}
