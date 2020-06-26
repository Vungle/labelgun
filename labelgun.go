package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: labelgun -stderrthreshold=[INFO|WARN|FATAL]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	// NOTE: This next line is key you have to call flag.Parse() for the command line
	// options or "flags" that are defined in the glog module to be picked up.
	flag.Parse()
}

func interval() int64 {
	val, _ := strconv.ParseInt(os.Getenv("LABELGUN_INTERVAL"), 10, 64)
	if val == 0 {
		return 60
	}
	return val
}

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	for {
		// Get Kube Nodes
		clientset := kubeClient(config)
		nodes, err := clientset.CoreV1().Nodes().List(v1.ListOptions{})
		if err != nil {
			log.Fatalf(err.Error())
		}

		// Get EC2 metadata
		metadata := ec2metadata.New(session.New())

		region, err := metadata.Region()
		if err != nil {
			log.Fatalf("Unable to retrieve the region from the EC2 instance %v\n", err)
		}

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

		for _, node := range nodes.Items {

			// Here we create an input that will filter any instances that aren't either
			// of these two states. This is generally what we want
			params := &ec2.DescribeInstancesInput{
				Filters: []*ec2.Filter{
					&ec2.Filter{
						Name: aws.String("private-dns-name"),
						Values: []*string{
							&node.Name,
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
				break
			}

			// Apply EC2 Tags
			nodeName := node.Name
			inst := resp.Reservations[0].Instances[0]

			for _, keys := range inst.Tags {
				tagKey := tagToLabel(*keys.Key)
				tagValue := tagToLabel(*keys.Value)

				if tagKey == "" || tagValue == "" {
					continue
				}

				label(nodeName, tagKey, tagValue)
			}
		}
		// Sleep until interval
		next(interval())
	}
}

func tagToLabel(item string) string {
	parsed, err := strconv.Unquote(string(awsutil.Prettify(item)))
	if err != nil {
		log.Error(err)
		return ""
	}
	parsed = strings.Replace(parsed, ":", ".", -1)
	if len(parsed) > 63 {
		return ""
	}
	return parsed
}

func next(interval int64) {
	log.Infof("Sleeping for %d seconds", interval)
	time.Sleep(time.Duration(interval) * time.Second)
}

func label(nodeName string, tagKey string, tagValue string) {
	log.Infoln(fmt.Sprintf("kubectl node %s %s=%s", nodeName, tagKey, tagValue))

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	clientset := kubeClient(config)
	node, err := clientset.CoreV1().Nodes().Get(nodeName)
	if err != nil {
		log.Fatalf(err.Error())
	}

	labels := node.GetLabels()
	labels[tagKey] = tagValue

	_, err = clientset.CoreV1().Nodes().Update(node)
	if err != nil {
		log.Error(err)
	}
}

func kubeClient(config *rest.Config) *kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return clientset
}
