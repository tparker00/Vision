package lookup

import (
	"fmt"

	"github.com/miekg/dns"
)

//ServiceLookup - Lookup srv records in DNS
func ServiceLookup(service string, dnsServer string) ([]string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(service+".", dns.TypeSRV)
	m.RecursionDesired = true
	r, _, err := c.Exchange(m, dnsServer+":53")

	if err != nil {
		fmt.Printf("Error happened: %s", err)
		return nil, err
	}
	var records []string

	for _, a := range r.Answer {
		records = append(records, a.(*dns.SRV).Target)
	}
	return records, nil
}
