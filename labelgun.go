package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codeskyblue/go-sh"
	"os"
	"strconv"
	"strings"
	"time"
)

var kube_master string
var node string
var svc = ec2.New(session.New(), &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
var meta *ec2.Instance

func main() {
	interval, _ := strconv.ParseInt(os.Getenv("LABELGUN_INTERVAL"), 10, 64)
	kube_master = os.Getenv("KUBE_MASTER")
	for {
		// Get Kube Node name
		n, _ := sh.Command("kubectl", "-s", kube_master, "describe", "pod", os.Getenv("HOSTNAME")).Command("grep", "Node").Command("awk", "{print $2}").Command("sed", "s@/.*@@").Output()
		node = string(n)
		node = strings.TrimSpace(node)
		fmt.Println(node)

		// Get instance id
		instance_id, _ := sh.Command("curl", "-s", "http://169.254.169.254/latest/meta-data/instance-id").Output()
		fmt.Println(string(instance_id))

		// Get AWS instance metadata
		params := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{
				aws.String(string(instance_id)),
			},
		}
		resp, err := svc.DescribeInstances(params)

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return
		}

		// Pretty-print the response data.
		//fmt.Println(resp.Reservations[0].Instances[0].InstanceType)
		meta = resp.Reservations[0].Instances[0]

		// Apply Availability Zone
		availabilityZone, _ := strconv.Unquote(string(awsutil.Prettify(meta.Placement.AvailabilityZone)))
		label("AvailabilityZone=" + availabilityZone)

		// Apply Instance Type
		instanceType, _ := strconv.Unquote(string(awsutil.Prettify(meta.InstanceType)))
		label("InstanceType=" + instanceType)

		// Apply EC2 Tags
		tags := meta.Tags
		for _, tag := range tags {
			label(*tag.Key + "=" + *tag.Value)
		}
		// Sleep until interval
		fmt.Println("Sleeping for " + os.Getenv("LABELGUN_INTERVAL") + " seconds")
		time.Sleep(time.Duration(interval) * time.Second)
	}

}

/*
func get_cores() string {
	cores, _ := sh.Command("nproc").Output()
	return "Cores=" + string(cores)
}
*/

func label(label_name string) {
	session := sh.NewSession()
	session.ShowCMD = true
	session.Command("kubectl", "-s", kube_master, "label", "node", node, label_name, "--overwrite").Run()
}
