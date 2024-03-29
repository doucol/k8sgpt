package calico

import (
	"github.com/fatih/color"
	"github.com/k8sgpt-ai/k8sgpt/pkg/common"
	"github.com/k8sgpt-ai/k8sgpt/pkg/kubernetes"
	"github.com/spf13/viper"
	operatorv1 "github.com/tigera/operator/api/v1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type StatusAnalyzer struct{}

func NewStatusAnalyzer() *StatusAnalyzer {
	return &StatusAnalyzer{}
}

func (s *StatusAnalyzer) Analyze(a common.Analyzer) ([]common.Result, error) {
	ret := []common.Result{}
	kubecontext := viper.GetString("kubecontext")
	kubeconfig := viper.GetString("kubeconfig")
	kc, err := kubernetes.NewClient(kubecontext, kubeconfig)
	if err != nil {
		color.Red("Error initialising kubernetes client config: %v", err)
		os.Exit(1)
	}
	ctrl := kc.GetCtrlClient()
	//ctrl := a.Client.GetCtrlClient()
	if err := operatorv1.AddToScheme(ctrl.Scheme()); err != nil {
		color.Red("Error initialising kubernetes client config: %v", err)
		os.Exit(1)
	}

	tigeraStatusList := operatorv1.TigeraStatusList{}
	err = ctrl.List(a.Context, &tigeraStatusList, &client.ListOptions{})
	if err != nil {
		return ret, err
	}

	for _, ts := range tigeraStatusList.Items {
		for _, cond := range ts.Status.Conditions {
			//if (cond.Type == operatorv1.ComponentReady || cond.Type == operatorv1.ComponentAvailable) && cond.Status != operatorv1.ConditionTrue {
			ret = append(ret, common.Result{
				Name: ts.Name,
				Kind: ts.Kind,
				Error: []common.Failure{
					{
						Text:          cond.Message,
						KubernetesDoc: "",
						Sensitive:     nil,
					},
				},
				Details:      cond.Message,
				ParentObject: "",
			})
			//}
		}
	}
	return ret, nil
}
