package e2e

import (
	"flag"
	//"fmt"
	"os"
	"path/filepath"
	"testing"
	//framework "github.com/operator-framework/operator-sdk/pkg/test"
	//"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	//"golang.org/x/net/context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestExample(t *testing.T) {

	require.True(t, true)

	t.Run("subtest example", func(t *testing.T) {
		require.True(t, true)

		// how to pick up configuration from jenkins
		// username, token and url, create namespace and delete component, needs to make calls to same api with same token
		//
		// how to read configurations from jenkins -- ask QE PEOPLE how to do this, sever url needs to be a configuration variable as well
		//
		// setup and install http://operator-hub-shbose-preview1-stage.b542.starter-us-east-2a.openshiftapps.com/operator/devopsconsole.v0.1.0 kubectl create -f to apply -f
		//operator registry needs to talk to the website above. - every push to master, should come and update some yaml file to say a newer version is available, then you should see the new version in openshift.  needs to be run as a CD job later

		// if you pass environment variables use them with openshift4 or else use the kubeconfig with minishift
		var kubeconfig *string

		if flag.Lookup("kubeconfig") != nil {
			tmp := flag.Lookup("kubeconfig").Value.(flag.Getter).Get().(string)
			kubeconfig = &tmp
		} else {
			if home := homeDir(); home != "" {
				kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}
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

		newNamespaceName := "test-namespace"
		_, err = clientset.CoreV1().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: newNamespaceName}})
		if err != nil {
			panic(err.Error())
		}

		createdNamespace, err := clientset.CoreV1().Namespaces().Get(newNamespaceName, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}

		createdNamespaceName := createdNamespace.GetName()
		require.Equal(t, newNamespaceName, createdNamespaceName)

		err = clientset.CoreV1().Namespaces().Delete(createdNamespaceName, &metav1.DeleteOptions{})
		if err != nil {
			panic(err.Error())
		}
	})
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
