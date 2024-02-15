package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// loops through network devices and tells us what one is our current dns address
func main() {

	// gather the fqdn
	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	fqdn := strings.ToLower(fmt.Sprintf("%v.%v", host, os.Getenv("USERDNSDOMAIN")))

	// ask dns for our IP
	dns, err := net.LookupIP(fqdn)
	if err != nil {
		log.Fatal(err)
	}

	// what devices match?
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, ifc := range interfaces {
		addresses, err := ifc.Addrs()
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		for _, addr := range addresses {
			var ip net.IP
			switch a := addr.(type) {
			case *net.IPNet:
				ip = a.IP
			case *net.IPAddr:
				ip = a.IP
			}

			if ip.To4() != nil {
				for _, dnsAddr := range dns {
					if ip.String() == dnsAddr.String() {
						fmt.Printf("\tDNS device: %v (%v)\n", ifc.Name, ip.String())
					}
				}
			}
		}
	}

}
