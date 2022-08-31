package main

import (
	"encoding/json"
	"fmt"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"os"
	"strings"
)

var (
	GroupName = os.Getenv("GROUP_NAME")
	regru     = RegruClient{os.Getenv("REGRU_USERNAME"), os.Getenv("REGRU_PASSWORD"), os.Getenv("REGRU_ZONE")}
)

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	cmd.RunWebhookServer(GroupName,
		&regruDNSProviderSolver{},
	)

}

type regruDNSProviderSolver struct {
	client *kubernetes.Clientset
}

type regruDNSProviderConfig struct {
	RegruAPIPasswordSecretRef cmmeta.SecretKeySelector `json:"regruPasswordSecretRef"`
}

func (c *regruDNSProviderSolver) Name() string {
	return "regru-dns"
}

func (c *regruDNSProviderSolver) Present(challengeRequest *v1alpha1.ChallengeRequest) error {
	klog.Infof("call function Present: namespace=%s, zone=%s, fqdn=%s", challengeRequest.ResourceNamespace, challengeRequest.ResolvedZone, challengeRequest.ResolvedFQDN)
	cfg, err := loadConfig(challengeRequest.Config)
	if err != nil {
		return fmt.Errorf("unable to load config: %v", err)
	}

	klog.Infof("decoded configuration %v", cfg)

	regruClient := NewRegruCient(regru.username, regru.password, regru.zone)

	entry, domain := c.getDomainAndEntry(challengeRequest)
	klog.Infof("present for entry=%s, domain=%s", entry, domain)

	err = regruClient.createTXT(domain, challengeRequest.Key)
	if err != nil {
		return fmt.Errorf("unable to check TXT record: %v", err)
	}
	return nil
}

func (c *regruDNSProviderSolver) CleanUp(challengeRequest *v1alpha1.ChallengeRequest) error {
	klog.Infof("call function CleanUp: namespace=%s, zone=%s, fqdn=%s",
		challengeRequest.ResourceNamespace, challengeRequest.ResolvedZone, challengeRequest.ResolvedFQDN)
	cfg, err := loadConfig(challengeRequest.Config)
	if err != nil {
		return fmt.Errorf("unable to load config: %v", err)
	}

	klog.Infof("decoded configuration %v", cfg)

	regruClient := NewRegruCient(regru.username, regru.password, regru.zone)
	entry, domain := c.getDomainAndEntry(challengeRequest)
	klog.Infof("present for entry=%s, domain=%s", entry, domain)

	err = regruClient.deleteTXT(domain, challengeRequest.Key)
	if err != nil {
		return fmt.Errorf("unable to check TXT record: %v", err)
	}
	return nil

}

func (c *regruDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, _ <-chan struct{}) error {
	klog.Infof("call function Initialize")
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return fmt.Errorf("unable to get k8s client: %v", err)
	}
	c.client = cl
	return nil
}

func (c *regruDNSProviderSolver) getDomainAndEntry(challengeRequest *v1alpha1.ChallengeRequest) (string, string) {
	entry := strings.TrimSuffix(challengeRequest.ResolvedFQDN, challengeRequest.ResolvedZone)
	entry = strings.TrimSuffix(entry, ".")
	domain := strings.TrimSuffix(challengeRequest.ResolvedZone, ".")
	return entry, domain

}

func loadConfig(cfgJSON *extapi.JSON) (regruDNSProviderConfig, error) {
	cfg := regruDNSProviderConfig{}
	if cfgJSON == nil {
		return cfg, nil
	}

	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		klog.Errorf("error decoding solver config: %v", err)
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}
	return cfg, nil
}
