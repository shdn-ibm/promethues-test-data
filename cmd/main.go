package main

import (
	"os"
	"os/signal"
	"syscall"

	"k8s.io/klog"

	"github.com/shdn-ibm/promethues-test-data/pkg/prome"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)
 
var scheme = runtime.NewScheme()
 
func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}
 
func main() {
	klog.InitFlags(nil)
 
	var err = false

	if err {
		goto error_out
	}
	go prome.RunExporter()
 
error_out:
	sigs := make(chan os.Signal)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	klog.Info("Awaiting signal to exit")
	go func() {
		sig := <-sigs
		klog.Infof("Received signal: %+v, clean up...", sig)
		done <- true
	}()

	// exiting
	<-done
	klog.Info("Exiting")
}
 