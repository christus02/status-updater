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
	"os"

	"k8s.io/client-go/rest"
)

var LOADBALANCER_IP_ANNOTATION = GetEnvWithFallback("LOADBALANCER_IP_ANNOTATION", "status.service.com/loadbalancer-ip")
var LOADBALANCER_HOSTNAME_ANNOTATION = GetEnvWithFallback("LOADBALANCER_HOSTNAME_ANNOTATION", "status.service.com/loadbalancer-hostname")

var clientset *kubernetes.Clientset

func GetEnvWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

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
	service_name := svc.ObjectMeta.Name
	service_namespace := svc.ObjectMeta.Namespace
	loadbalancerIPValue, ipOk := svc.ObjectMeta.Annotations[LOADBALANCER_IP_ANNOTATION]
	loadbalancerHostnameValue, hostnameOk := svc.ObjectMeta.Annotations[LOADBALANCER_HOSTNAME_ANNOTATION]
	copiedSvc := svc.DeepCopy()
	if hostnameOk && ipOk {
		klog.Infof("Service: %s, Namespace: %s, Hostname : %s, IP: %s", service_name, service_namespace, loadbalancerHostnameValue, loadbalancerIPValue)
		copiedSvc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{Hostname: loadbalancerHostnameValue}, {IP: loadbalancerIPValue}}
		res, err := clientset.CoreV1().Services(service_namespace).UpdateStatus(context.TODO(), copiedSvc, metav1.UpdateOptions{FieldManager: "christus-controller"})
		if err != nil {
			klog.Error(err)
		} else {
			klog.Info("Successfully Updated LoadBalancer Status")
			klog.Info("Response Dump: ", res)
		}
	} else if hostnameOk {
		klog.Infof("Service: %s, Namespace: %s, Hostname : %s", service_name, service_namespace, loadbalancerHostnameValue)
		copiedSvc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{Hostname: loadbalancerHostnameValue}}
		res, err := clientset.CoreV1().Services(service_namespace).UpdateStatus(context.TODO(), copiedSvc, metav1.UpdateOptions{FieldManager: "christus-controller"})
		if err != nil {
			klog.Error(err)
		} else {
			klog.Info("Successfully Updated LoadBalancer Status")
			klog.Info("Response Dump: ", res)
		}
	} else if ipOk {
		klog.Infof("Service: %s, Namespace: %s, IP : %s", service_name, service_namespace, loadbalancerIPValue)
		copiedSvc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{IP: loadbalancerIPValue}}
		res, err := clientset.CoreV1().Services(service_namespace).UpdateStatus(context.TODO(), copiedSvc, metav1.UpdateOptions{FieldManager: "christus-controller"})
		if err != nil {
			klog.Error(err)
		} else {
			klog.Info("Successfully Updated LoadBalancer Status")
			klog.Info("Response Dump: ", res)
		}
	}

}

func onUpdate(obj interface{}, newObj interface{}) {
	newSvc := newObj.(*corev1.Service)
	service_name := newSvc.ObjectMeta.Name
	service_namespace := newSvc.ObjectMeta.Namespace
	loadbalancerIPValue, ipOk := newSvc.ObjectMeta.Annotations[LOADBALANCER_IP_ANNOTATION]
	loadbalancerHostnameValue, hostnameOk := newSvc.ObjectMeta.Annotations[LOADBALANCER_HOSTNAME_ANNOTATION]
	copiedSvc := newSvc.DeepCopy()
	if hostnameOk && ipOk {
		klog.Infof("Service: %s, Namespace: %s, Hostname : %s, IP: %s", service_name, service_namespace, loadbalancerHostnameValue, loadbalancerIPValue)
		copiedSvc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{Hostname: loadbalancerHostnameValue}, {IP: loadbalancerIPValue}}
		res, err := clientset.CoreV1().Services(service_namespace).UpdateStatus(context.TODO(), copiedSvc, metav1.UpdateOptions{FieldManager: "christus-controller"})
		if err != nil {
			klog.Error(err)
		} else {
			klog.Info("Successfully Updated LoadBalancer Status")
			klog.Info("Response Dump: ", res)
		}
	} else if hostnameOk {
		klog.Infof("Service: %s, Namespace: %s, Hostname : %s", service_name, service_namespace, loadbalancerHostnameValue)
		copiedSvc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{Hostname: loadbalancerHostnameValue}}
		res, err := clientset.CoreV1().Services(service_namespace).UpdateStatus(context.TODO(), copiedSvc, metav1.UpdateOptions{FieldManager: "christus-controller"})
		if err != nil {
			klog.Error(err)
		} else {
			klog.Info("Successfully Updated LoadBalancer Status")
			klog.Info("Response Dump: ", res)
		}
	} else if ipOk {
		klog.Infof("Service: %s, Namespace: %s, IP : %s", service_name, service_namespace, loadbalancerIPValue)
		copiedSvc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{IP: loadbalancerIPValue}}
		res, err := clientset.CoreV1().Services(service_namespace).UpdateStatus(context.TODO(), copiedSvc, metav1.UpdateOptions{FieldManager: "christus-controller"})
		if err != nil {
			klog.Error(err)
		} else {
			klog.Info("Successfully Updated LoadBalancer Status")
			klog.Info("Response Dump: ", res)
		}
	}
}

func main() {
	klog.Info("***** Status Updater Started *****")
	InitWatcher()
}
