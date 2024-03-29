package calico

import (
	"context"
	"errors"
	"github.com/fatih/color"
	"github.com/k8sgpt-ai/k8sgpt/pkg/common"
	"github.com/k8sgpt-ai/k8sgpt/pkg/kubernetes"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

const (
	StatusAnalyzerName = "CalicoStatus"
)

type Calico struct{}

func NewCalico() *Calico {
	return &Calico{}
}

func (c *Calico) Deploy(_ string) error {
	color.Green("Activating calico integration...")
	client := getKubernetesClient().GetClient()
	ctx := context.Background()
	_, err := client.AppsV1().DaemonSets("calico-system").Get(ctx, "calico-node", metav1.GetOptions{})
	if err != nil {
		color.Yellow(`Calico installation not found. Please ensure Calico is deployed to analyze.`)
		return errors.New("no calico installation found")
	}
	color.Green("Found existing Calico installation")
	return nil
}

func (c *Calico) UnDeploy(namespace string) error {
	color.Yellow("Integration will leave Calico resources deployed.")
	return nil
}

func (c *Calico) AddAnalyzer(mergedMap *map[string]common.IAnalyzer) {
	(*mergedMap)[StatusAnalyzerName] = NewStatusAnalyzer()
}

func (c *Calico) GetAnalyzerName() []string {
	return []string{
		StatusAnalyzerName,
	}
}

func (c *Calico) GetNamespace() (string, error) {
	return "", nil
}

func (c *Calico) OwnsAnalyzer(s string) bool {
	for _, an := range c.GetAnalyzerName() {
		if s == an {
			return true
		}
	}
	return false
}

func (c *Calico) IsActivate() bool {
	filters := viper.GetStringSlice("active_filters")
	for _, filter := range filters {
		for _, analyzer := range c.GetAnalyzerName() {
			if filter == analyzer {
				return true
			}
		}
	}
	return false
}

func getKubernetesClient() *kubernetes.Client {
	kubecontext := viper.GetString("kubecontext")
	kubeconfig := viper.GetString("kubeconfig")
	k, err := kubernetes.NewClient(kubecontext, kubeconfig)
	if err != nil {
		color.Red("Error initialising kubernetes client: %v", err)
		os.Exit(1)
	}
	return k
}
