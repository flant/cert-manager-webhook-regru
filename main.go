package main

import (
	"fmt"
	"os"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

const providerName = "regru-dns"

var (
	GroupName = os.Getenv("GROUP_NAME")
	regru     = RegruClient{os.Getenv("REGRU_USERNAME"), os.Getenv("REGRU_PASSWORD"), ""}
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
	return providerName
}

func (c *regruDNSProviderSolver) Present(challengeRequest *v1alpha1.ChallengeRequest) error {
	klog.Infof("Call function Present: namespace=%s, zone=%s, fqdn=%s", challengeRequest.ResourceNamespace, challengeRequest.ResolvedZone, challengeRequest.ResolvedFQDN)
	//_, err := loadConfig(challengeRequest.Config)
	//if err != nil {
	//	return fmt.Errorf("unable to load config: %v", err)
	//}
	//
	//klog.Infof("decoded configuration %v", cfg)

	zone, err := getDomainFromZone(challengeRequest.ResolvedZone, challengeRequest.ResolvedFQDN)
	if err != nil {
		return fmt.Errorf("unable to initialize reg.ru client, because unable to get root zone from domains: %w", err)
	}

	klog.Infof("Using reg.ru client with username %s and zone %s", regru.username, zone)

	regruClient := NewRegruClient(regru.username, regru.password, zone)

	klog.Infof("present for entry=%s, domain=%s, key=%s", challengeRequest.ResolvedFQDN, zone, challengeRequest.Key)
	if err := regruClient.createTXT(challengeRequest.ResolvedFQDN, challengeRequest.Key); err != nil {
		return fmt.Errorf("unable to create TXT record: %v", err)
	}

	return nil
}

func (c *regruDNSProviderSolver) CleanUp(challengeRequest *v1alpha1.ChallengeRequest) error {
	klog.Infof("Call function CleanUp: namespace=%s, zone=%s, fqdn=%s",
		challengeRequest.ResourceNamespace, challengeRequest.ResolvedZone, challengeRequest.ResolvedFQDN)
	//cfg, err := loadConfig(challengeRequest.Config)
	//if err != nil {
	//	return fmt.Errorf("unable to load config: %v", err)
	//}
	//
	//klog.Infof("decoded configuration %v", cfg)

	zone, err := getDomainFromZone(challengeRequest.ResolvedZone, challengeRequest.ResolvedFQDN)
	if err != nil {
		return fmt.Errorf("unable to initialize reg.ru client, because unable to get root zone from domains: %w", err)
	}

	klog.Infof("Using reg.ru client with username %s and zone %s", regru.username, zone)

	regruClient := NewRegruClient(regru.username, regru.password, zone)
	klog.Infof("delete entry=%s, domain=%s, key=%s", challengeRequest.ResolvedFQDN, zone, challengeRequest.Key)

	if err := regruClient.deleteTXT(challengeRequest.ResolvedFQDN, challengeRequest.Key); err != nil {
		return fmt.Errorf("unable to delete TXT record: %v", err)
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

//func loadConfig(cfgJSON *extapi.JSON) (regruDNSProviderConfig, error) {
//	cfg := regruDNSProviderConfig{}
//	if cfgJSON == nil {
//		return cfg, nil
//	}
//
//	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
//		klog.Errorf("error decoding solver config: %v", err)
//		return cfg, fmt.Errorf("error decoding solver config: %v", err)
//	}
//	return cfg, nil
//}
