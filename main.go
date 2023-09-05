package main

import (
	"fmt"
	"os"
	"strings"

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

	regruClient := NewRegruClient(regru.username, regru.password, getDomainFromZone(challengeRequest.ResolvedZone))

	klog.Infof("present for entry=%s, domain=%s, key=%s", challengeRequest.ResolvedFQDN, getDomainFromZone(challengeRequest.ResolvedZone), challengeRequest.Key)

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

	regruClient := NewRegruClient(regru.username, regru.password, getDomainFromZone(challengeRequest.ResolvedZone))
	klog.Infof("delete entry=%s, domain=%s, key=%s", challengeRequest.ResolvedFQDN, getDomainFromZone(challengeRequest.ResolvedZone), challengeRequest.Key)

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

// getDomainFromZone returns second-level domain name from ResolvedZone without last dot.
// reg.ru api requires to specify the second-level domain in the request
func getDomainFromZone(zone string) string {
	parts := strings.Split(zone[0:len(zone)-1], ".")
	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}
