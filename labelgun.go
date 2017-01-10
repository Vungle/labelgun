package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codeskyblue/go-sh"
	log "github.com/golang/glog"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"fmt"
	"bytes"
)

var node string

func usage() {
	fmt.Fprintf(os.Stderr, "usage: labelgun -stderrthreshold=[INFO|WARN|FATAL]\n", )
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	// NOTE: This next line is key you have to call flag.Parse() for the command line
	// options or "flags" that are defined in the glog module to be picked up.
	flag.Parse()
}

func errThreshold() string {
	errThreshold := os.Getenv("LABELGUN_ERR_THRESHOLD")
	if errThreshold == "" {
		return "INFO"
	}
	return errThreshold
}

func interval() int64 {
	val, _ := strconv.ParseInt(os.Getenv("LABELGUN_INTERVAL"), 10, 64)
	if val == 0 {
		return 60
	}
	return val
}

func main() {
	for {
		// Get Kube Node name
		n, err := sh.Command("kubectl", "describe", "pod", os.Getenv("HOSTNAME")).Command("grep", "Node").Command("awk", "{print $2}").Command("sed", "s@/.*@@").Output()
		if err != nil {
			log.Fatal(err)
		}
		node = string(n)
		node = strings.TrimSpace(node)
		log.Infoln(node)

		// Get EC2 metadata
		metadata := ec2metadata.New(session.New())

		region, err := metadata.Region()
		if err != nil {
			log.Fatalf("Unable to retrieve the region from the EC2 instance %v\n", err)
		}

		doc, err := metadata.GetInstanceIdentityDocument()
		if err != nil {
			log.Fatalf("Unable to retrieve the metadata from the EC2 instance %v\n", err)
		}

		// Apply Availability Zone
		availabilityZone, _ := strconv.Unquote(string(awsutil.Prettify(doc.AvailabilityZone)))
		go label(node,"AvailabilityZone", availabilityZone)
		log.Infoln(availabilityZone)

		// Apply Instance Type
		instanceType, _ := strconv.Unquote(string(awsutil.Prettify(doc.InstanceType)))
		go label(node,"InstanceType", instanceType)
		log.Infoln(instanceType)

		creds := credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{},
				&ec2rolecreds.EC2RoleProvider{Client: metadata},
			})

		awsConfig := aws.NewConfig()
		awsConfig.WithCredentials(creds)
		awsConfig.WithRegion(region)
		sess, err := session.NewSession(awsConfig)
		if err != nil {
			log.Fatal(err)
		}

		svc := ec2.New(sess)

		// Here we create an input that will filter any instances that aren't either
		// of these two states. This is generally what we want
		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name: aws.String("instance-id"),
					Values: []*string{
						&doc.InstanceID,
					},
				},
			},
		}

		resp, _ := svc.DescribeInstances(params)
		if err != nil {
			log.Fatal(err)
		}
		if len(resp.Reservations) < 1 || len(resp.Reservations[0].Instances) < 1 {
			// Might due to "Request body type has been overwritten. May cause race conditions"
			next(interval())
			break;
		}

		// Apply EC2 Tags
		inst := resp.Reservations[0].Instances[0];
		for _, keys := range inst.Tags {
			tagKey, _ := strconv.Unquote(string(awsutil.Prettify(*keys.Key)))
			tagValue, _ := strconv.Unquote(string(awsutil.Prettify(*keys.Value)))
			go label(node, tagKey, tagValue)
		}

		// Sleep until interval
		next(interval())
	}
}

func next(interval int64) {
	log.Infof("Sleeping for %d seconds", interval)
	time.Sleep(time.Duration(interval) * time.Second)
}

func label(node string, label_key string, label_value string) {
	if node == "" {
		log.Fatalf("node is empty!")
	}

	session := sh.NewSession()
	if strings.EqualFold(errThreshold(), "INFO") {
		session.ShowCMD = true
	} else {
		var outbuf, errbuf bytes.Buffer
		session.Stdout = &outbuf
		session.Stderr = &errbuf
	}

	_, err := session.Command("kubectl", "label", "node", node, label_key+"="+label_value, "--overwrite").Output()
	if err != nil {
		log.Warningln("oops, something was too hard", err)
		return
	}
}
