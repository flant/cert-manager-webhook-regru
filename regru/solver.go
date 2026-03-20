package regru

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"golang.org/x/net/publicsuffix"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const solverName = "regru-dns"

// Solver implements the cert-manager DNS01 webhook solver for the reg.ru DNS provider.
type Solver struct {
	client *kubernetes.Clientset
}

// Name returns the solver name used to select this webhook in cert-manager challenges.
func (s *Solver) Name() string {
	return solverName
}

// Initialize sets up the Kubernetes clientset required for accessing cluster resources.
func (s *Solver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	slog.Default().Info("initializing webhook solver")

	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		slog.Default().Error("failed to create k8s clientset", "error", err)
		return fmt.Errorf("create k8s clientset: %w", err)
	}
	s.client = cl

	slog.Default().Info("webhook solver initialized")
	return nil
}

// Present creates the TXT DNS record required to fulfill an ACME DNS-01 challenge.
func (s *Solver) Present(ch *v1alpha1.ChallengeRequest) error {
	if ch == nil {
		slog.Default().Error("present: challenge request is nil")
		return fmt.Errorf("challenge request is nil")
	}
	start := time.Now()
	zoneDomain, err := extractDomain(ch.ResolvedZone)

	if err != nil {
		slog.Default().Error("present: cannot extract registrable domain from resolved zone",
			"error", err,
			"zone", ch.ResolvedZone,
			"dns_name", ch.DNSName,
			"fqdn", ch.ResolvedFQDN,
		)
		return fmt.Errorf("invalid resolved zone %q: cannot extract registrable domain: %w", ch.ResolvedZone, err)
	}
	if zoneDomain == "" {
		slog.Default().Error("present: cannot extract registrable domain from resolved zone",
			"zone", ch.ResolvedZone,
			"dns_name", ch.DNSName,
			"fqdn", ch.ResolvedFQDN,
		)
		return fmt.Errorf("invalid resolved zone %q: cannot extract registrable domain", ch.ResolvedZone)
	}

	slog.Default().Info("received ACME challenge, need to create TXT record",
		"dns_name", ch.DNSName,
		"zone", ch.ResolvedZone,
		"fqdn", ch.ResolvedFQDN,
		"zoneDomain", zoneDomain,
	)

	if err := regruClient.addTXTRecord(context.Background(), zoneDomain, ch.ResolvedFQDN, ch.Key); err != nil {
		slog.Default().Error("present: failed to create TXT record",
			"error", err,
			"zoneDomain", zoneDomain,
			"fqdn", ch.ResolvedFQDN,
		)
		return fmt.Errorf("create TXT record for domain %q fqdn %q: %w", zoneDomain, ch.ResolvedFQDN, err)
	}

	slog.Default().Info("ACME challenge TXT record created successfully",
		"zoneDomain", zoneDomain,
		"fqdn", ch.ResolvedFQDN,
		"duration_ms", time.Since(start).Milliseconds(),
	)
	return nil
}

// CleanUp removes the TXT DNS record after the ACME DNS-01 challenge has been verified.
func (s *Solver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	if ch == nil {
		slog.Default().Error("cleanup: challenge request is nil")
		return fmt.Errorf("challenge request is nil")
	}
	start := time.Now()
	zoneDomain, err := extractDomain(ch.ResolvedZone)

	if err != nil {
		slog.Default().Error("cleanup: cannot extract registrable domain from resolved zone",
			"error", err,
			"zone", ch.ResolvedZone,
			"dns_name", ch.DNSName,
			"fqdn", ch.ResolvedFQDN,
		)
		return fmt.Errorf("invalid resolved zone %q: cannot extract registrable domain: %w", ch.ResolvedZone, err)
	}
	if zoneDomain == "" {
		slog.Default().Error("cleanup: cannot extract registrable domain from resolved zone",
			"zone", ch.ResolvedZone,
			"dns_name", ch.DNSName,
			"fqdn", ch.ResolvedFQDN,
		)
		return fmt.Errorf("invalid resolved zone %q: cannot extract registrable domain", ch.ResolvedZone)
	}

	slog.Default().Info("received cleanup request, need to delete TXT record",
		"dns_name", ch.DNSName,
		"zone", ch.ResolvedZone,
		"fqdn", ch.ResolvedFQDN,
		"zoneDomain", zoneDomain,
	)

	if err := regruClient.deleteTXTRecord(context.Background(), zoneDomain, ch.ResolvedFQDN, ch.Key); err != nil {
		slog.Default().Error("cleanup: failed to delete TXT record",
			"error", err,
			"zoneDomain", zoneDomain,
			"fqdn", ch.ResolvedFQDN,
		)
		return fmt.Errorf("delete TXT record for domain %q fqdn %q: %w", zoneDomain, ch.ResolvedFQDN, err)
	}

	slog.Default().Info("ACME challenge TXT record deleted successfully",
		"fqdn", ch.ResolvedFQDN,
		"zoneDomain", zoneDomain,
		"duration_ms", time.Since(start).Milliseconds(),
	)
	return nil
}

// extractDomain extracts the registrable domain from a zone name.
// e.g., "example.com." -> "example.com"
// e.g., "sub.example.com." -> "example.com"
// Returns empty string with nil error for empty input or ".", and returns an error
// for invalid zones (e.g. single label) or publicsuffix evaluation failures.
func extractDomain(zone string) (string, error) {
	z := strings.TrimSuffix(zone, ".")
	if z == "" || z == "." {
		return "", nil
	}

	d, err := publicsuffix.EffectiveTLDPlusOne(z)
	if err != nil {
		return "", err
	}
	return d, nil
}
