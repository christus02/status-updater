package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
	//"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
	"os"
)

var SERVICE_NAME = os.Getenv("SERVICE_NAME")
var SERVICE_NAMESPACE = os.Getenv("SERVICE_NAMESPACE")
var EXTERNAL_ENDPOINT_TYPE_ANNOTATION = os.Getenv("EXTERNAL_ENDPOINT_TYPE_ANNOTATION")
var ENDPOINT_ANNOTATION = os.Getenv("ENDPOINT_ANNOTATION")
var ENDPOINT = ""
var ENDPOINT_TYPE = "hostname"

var clientset *kubernetes.Clientset

func InitWatcher() {

	config, err := rest.InClusterConfig() // Using in-cluster-config serviceAccount will be used
	//config, err := clientcmd.BuildConfigFromFlags("", "/home/someuser/.kube/config")
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
	klog.Info("Starting Service Informer - Watching for Service Events")
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
		endpointValue, ok := svc.ObjectMeta.Annotations[ENDPOINT_ANNOTATION]
		endpointTypeValue, typeOk := svc.ObjectMeta.Annotations[EXTERNAL_ENDPOINT_TYPE_ANNOTATION]
		if ok && typeOk {
			ENDPOINT = endpointValue
			ENDPOINT_TYPE = endpointTypeValue
			klog.Infof("Service: %s, Namespace: %s, ENDPOINT: %s, ENDPOINT TYPE: %s", SERVICE_NAME, SERVICE_NAMESPACE, ENDPOINT, ENDPOINT_TYPE)
			copiedSvc := svc.DeepCopy()
			if ENDPOINT_TYPE == "hostname" {
				klog.Info("Updating LoadBalancer Status for Endpoint Hostname: ", ENDPOINT)
				copiedSvc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{Hostname: ENDPOINT}}
			} else if ENDPOINT_TYPE == "ip" {
				klog.Info("Updating LoadBalancer Status for Endpoint IP: ", ENDPOINT)
				copiedSvc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{IP: ENDPOINT}}
			}
			res, err := clientset.CoreV1().Services(SERVICE_NAMESPACE).UpdateStatus(context.TODO(), copiedSvc, metav1.UpdateOptions{FieldManager: "christus-controller"})
			if err != nil {
				klog.Error(err)
			} else {
				klog.Info("Successfully Updated LoadBalancer Status")
				klog.Info("Response Dump: ", res)
			}

		} else {
			klog.Errorf("Service: %s, Namespace: %s, ENDPOINT or ENDPOINT_TYPE is *NOT SPECIFIED*", SERVICE_NAME, SERVICE_NAMESPACE)
		}

	}
}

func onUpdate(obj interface{}, newObj interface{}) {
	newSvc := newObj.(*corev1.Service)
	if newSvc.ObjectMeta.Name == SERVICE_NAME && newSvc.ObjectMeta.Namespace == SERVICE_NAMESPACE {
		endpointValue, ok := newSvc.ObjectMeta.Annotations[ENDPOINT_ANNOTATION]
		endpointTypeValue, typeOk := newSvc.ObjectMeta.Annotations[EXTERNAL_ENDPOINT_TYPE_ANNOTATION]
		if ok && typeOk {
			ENDPOINT = endpointValue
			ENDPOINT_TYPE = endpointTypeValue
			klog.Infof("Service: %s, Namespace: %s, ENDPOINT: %s, ENDPOINT TYPE: %s", SERVICE_NAME, SERVICE_NAMESPACE, ENDPOINT, ENDPOINT_TYPE)
			copiedSvc := newSvc.DeepCopy()
			if ENDPOINT_TYPE == "hostname" {
				klog.Info("Updating LoadBalancer Status for Endpoint Hostname: ", ENDPOINT)
				copiedSvc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{Hostname: ENDPOINT}}
			} else if ENDPOINT_TYPE == "ip" {
				klog.Info("Updating LoadBalancer Status for Endpoint IP: ", ENDPOINT)
				copiedSvc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{IP: ENDPOINT}}
			}
			res, err := clientset.CoreV1().Services(SERVICE_NAMESPACE).UpdateStatus(context.TODO(), copiedSvc, metav1.UpdateOptions{FieldManager: "christus-controller"})
			if err != nil {
				klog.Error(err)
			} else {
				klog.Info("Successfully Updated LoadBalancer Status")
				klog.Info("Response Dump: ", res)
			}

		} else {
			klog.Errorf("Service: %s, Namespace: %s, ENDPOINT or ENDPOINT_TYPE is *NOT SPECIFIED*", SERVICE_NAME, SERVICE_NAMESPACE)
		}

	}

}

func main() {
	InitWatcher()
}
