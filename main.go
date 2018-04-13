package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

var version = "0.1.4"

var (
	help      = flag.Bool("help", false, "Show this help.")
	service   = flag.String("service", "", "Governing service this pod is in.")
	discovery = flag.String("discovery", "kubedns", "Service Discovery this pod & service use. {kubedns|synapse}")
	toExclude = flag.String("exclude", "", "Excluded from seeds. Coma seperated list of pods")
	output    = flag.String("output", "stdout", "Write seed list in a file if specified.")
)

func main() {
	flag.Parse()

	var seeds []string
	var excluded []string
	separator := ","
	self, err := os.Hostname()
	if err != nil {
		fmt.Printf("Unable to find self hostname, this is mandatory to exclude it from the list")
		os.Exit(1)
	}

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
	excluded = append(excluded, self)

	switch *discovery {
	case "kubedns":
		cname, srvRecords, err := net.LookupSRV("", "", *service)
		if err != nil {
			fmt.Printf("Unable to do SRV query for %q: %v", *service, err)
			os.Exit(1)
		}
		for _, srvRecord := range srvRecords {
			if !contains(excluded, srvRecord.Target, cname) && !contains(excluded, self, cname) {
				seeds = append(seeds, srvRecord.Target)
			}
		}

		switch *output {
		case "stdout":
			fmt.Printf("%s", strings.Join(seeds, separator))
		default:
			if err := ioutil.WriteFile(*output, []byte(strings.Join(seeds, separator)), 0666); err != nil {
				fmt.Printf("Unable to write %s file: %v", *output, err)
				os.Exit(1)
			}
		}

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
