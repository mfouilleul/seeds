package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var version = "0.1.0"

var (
	help      = flag.Bool("help", false, "Show this help.")
	service   = flag.String("service", "", "Governing service this pod is in. (Required)")
	discovery = flag.String("discovery", "kubedns", "Service Discovery this pod & service use. {kubedns|synapse}")
	toExclude = flag.String("exclude", "", "Excluded from seeds. Coma seperated list of pods")
)

func main() {
	flag.Parse()

	var seeds []string
	var excluded []string
	separator := ","

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *service == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *toExclude != "" {
		excluded = strings.Split(*toExclude, ",")
	}

	switch *discovery {
	case "kubedns":
		cname, srvRecords, err := net.LookupSRV("", "", *service)
		if err != nil {
			fmt.Printf("Unable to do SRV query for %q: %v", *service, err)
			os.Exit(1)
		}
		for _, srvRecord := range srvRecords {
			if !contains(excluded, srvRecord.Target, cname) {
				seeds = append(seeds, srvRecord.Target)
			}
		}

		fmt.Printf("%s", strings.Join(seeds, separator))

	case "synapse": // TODO
		fmt.Printf("%v is not implemented yet\n", *discovery)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func contains(podList []string, searchPod string, cname string) bool {
	for _, value := range podList {
		value = fmt.Sprintf("%s.%s", value, cname)
		if value == searchPod {
			return true
		}
	}
	return false
}
