package elbv2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func statusLoadBalancerState(conn *elbv2.ELBV2, arn string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &elbv2.DescribeLoadBalancersInput{
			LoadBalancerArns: []*string{aws.String(arn)},
		}

		output, err := conn.DescribeLoadBalancers(input)

		if tfawserr.ErrCodeEquals(err, elbv2.ErrCodeLoadBalancerNotFoundException) {
			return nil, "", nil
		}
		if err != nil {
			return nil, "", err
		}

		if len(output.LoadBalancers) != 1 {
			return nil, "", fmt.Errorf("No load balancers found for %s", arn)
		}
		lb := output.LoadBalancers[0]

		return output, aws.StringValue(lb.State.Code), nil
	}
}

func statusTargetGroupState(conn *elbv2.ELBV2, arn string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &elbv2.DescribeTargetHealthInput{
			TargetGroupArn: aws.String(arn),
		}

		output, err := conn.DescribeTargetHealth(input)

		if err != nil {
			return nil, "", err
		}

		if len(output.TargetHealthDescriptions) != 1 {
			return nil, "", fmt.Errorf("No Target Group found for %s", arn)
		}
		tg := output.TargetHealthDescriptions[0]

		return output, aws.StringValue(tg.TargetHealth.State), nil

	}
}
