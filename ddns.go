package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"strings"
	"time"
	"github.com/kardianos/osext"
)

const (
	INTERVAL = 5 //Minute
)

var (
	config Settings
	configPath = flag.String("c", "conf/app.yml", "Specify a configuration file")
	showHelp = flag.Bool("h", false, "Show help")
)

func main() {
	flag.Parse()
	if *showHelp {
		flag.Usage()
		return
	}

	if !strings.HasPrefix(*configPath, "/") {
		pwd, _ := osext.ExecutableFolder()
		*configPath = fmt.Sprintf("%s%s%s", pwd, "/", *configPath)
	}

	print(*configPath)

	if err := Load(*configPath, &config); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	log.Println(config)

	dnsLoop()
}

func dnsLoop() {
	for {
		currentIP, err := getCurrentIP(config.IPService)

		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 3)
			continue
		}

		subDomainID, ip := getSubDomain(config.Domain, config.SubDomain, config.Type)

		if subDomainID == "" || ip == "" {
			fmt.Println("sub_domain:", subDomainID, ip)
			continue
		}

		if len(ip) > 0 && !strings.Contains(currentIP, ip) {
			updateIP(config.Domain, subDomainID, config.SubDomain, currentIP, config.Ttl)
		} else {
			fmt.Println("Current IP is same as domain IP, no need to update...")
		}

		//Interval is 5 minutes
		time.Sleep(time.Minute * INTERVAL)
	}
}