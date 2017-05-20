package main

import (
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

var PREFIX = "KRKN"

func main() {
	usage := `kraken - EC2/EBS utility

Usage:
  kraken [-l=<log-level>] [--quiet] [--dry]
  kraken -h | --help
  kraken --version

Options:
  --quiet              Suppress logs
  --log-level=<level>  Log level (debug, info, warn, error, fatal)
  --dry                Dry run
  -h --help            Show this screen.
  --version            Show version.`
	arguments, _ := docopt.Parse(usage, nil, true, "kraken 0.1", false)

	if arguments["--quiet"].(bool) {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}

	logLevel, ok := arguments["--log-level"].(string)
	if !ok {
		log.SetLevel(log.WarnLevel)
	} else {
		if level, err := log.ParseLevel(logLevel); err != nil {
			log.Fatalf("%v", err)
		} else {
			log.SetLevel(level)
		}
	}

	if arguments["--verbose"].(bool) {
	}

	log.SetFormatter(&log.TextFormatter{})

	dryRun := arguments["--dry"].(bool)

	log.Printf("%s", DrawAsciiBanner("WELCOME TO KRAKEN"))

	log.Printf("%s", DrawAsciiBanner("1: COLLECTING EC2 INFO"))
	ec2Instance := GetEc2InstanceData()

	log.Printf("%s", DrawAsciiBanner("2: COLLECTING EBS INFO"))
	ebsVolumes := MapEbsVolumes(&ec2Instance)

	log.Printf("%s", DrawAsciiBanner("3: ATTACHING EBS VOLS"))
	ebsVolumes = AttachEbsVolumes(ec2Instance, ebsVolumes, dryRun)

	log.Printf("%s", DrawAsciiBanner("4: MOUNTING ATTACHED VOLS"))

	if len(ebsVolumes) == 0 {
		log.Warn("Empty vols, nothing to do")
		os.Exit(0)
	}

	for volName, vols := range ebsVolumes {
		PrepAndMountDrives(volName, vols, dryRun)
	}
}
