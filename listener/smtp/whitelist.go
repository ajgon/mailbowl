package smtp

import (
	"net"

	"github.com/Masterminds/log-go"
)

type Whitelist struct {
	CIDRs []string
}

func NewWhitelist(cidrs []string) *Whitelist {
	validCIDRs := make([]string, 0)

	for _, cidr := range cidrs {
		_, _, err := net.ParseCIDR(cidr)

		if err != nil {
			log.Debugf("invalid smtp.whitelist.cidrs entry `%s`, removing", cidr)
		} else {
			validCIDRs = append(validCIDRs, cidr)
		}
	}

	return &Whitelist{
		CIDRs: validCIDRs,
	}
}
