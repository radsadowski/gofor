package validation

import (
	"net"

	"github.com/asaskevich/govalidator"
)

func IsValidIPAddress(address string) bool {
	result := net.ParseIP(address)
	return result != nil
}

func IsValidHostname(name string) bool {
	return govalidator.IsDNSName(name)
}
