package main

import (
	"errors"
	"fmt"
	"strings"
)

// getDomainFromZone returns second-level domain name from ResolvedZone without last dot.
// reg.ru api requires to specify the second-level domain in the request
func getDomainFromZone(domains ...string) (string, error) {
	for _, domain := range domains {
		parts := strings.Split(domain[0:len(domain)-1], ".")
		if len(parts) > 1 {
			return parts[len(parts)-2] + "." + parts[len(parts)-1], nil
		}
	}
	return "", errors.New(fmt.Sprintf("not enouth parts in domains to find root zone: %v", domains))
}
