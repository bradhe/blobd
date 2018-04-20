package aws

import (
	baseaws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

const DefaultAWSRegion = "us-east-1"

var (
	defaultSession        = session.Must(session.NewSession(baseaws.NewConfig()))
	defaultMetadataClient = ec2metadata.New(defaultSession)

	defaultCredentials = credentials.NewChainCredentials(
		[]credentials.Provider{
			&ec2rolecreds.EC2RoleProvider{
				Client: defaultMetadataClient,
			},
			&credentials.EnvProvider{},
		})

	defaultConfig = baseaws.NewConfig().
			WithRegion(getAWSRegion()).
			WithMaxRetries(3).
			WithCredentials(defaultCredentials)
)

func getAWSRegion() string {
	if region, err := defaultMetadataClient.Region(); err == nil {
		return DefaultAWSRegion
	} else {
		return region
	}
}

func getAWSSession() *session.Session {
	return session.Must(session.NewSession(defaultConfig))
}
