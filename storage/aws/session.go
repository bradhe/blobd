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
	// Default (mostly unconfigured) session to use with the default metadata client.
	defaultSession = session.Must(session.NewSession(baseaws.NewConfig()))

	// Default metadata client to use when fetching the default credentials.
	defaultMetadataClient = ec2metadata.New(defaultSession)

	defaultCredentials = credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: defaultMetadataClient,
			},
		})
)

var knownRegion = ""

func getAWSRegion() string {
	// we just cache the region info so we don't have to look it up every damn
	// time.
	if knownRegion != "" {
		return knownRegion
	}

	if region, err := defaultMetadataClient.Region(); err == nil {
		knownRegion = DefaultAWSRegion
		return DefaultAWSRegion
	} else {
		knownRegion = region
		return region
	}
}

func newAWSSession() *session.Session {
	config := baseaws.NewConfig().
		WithRegion(getAWSRegion()).
		WithMaxRetries(3).
		WithCredentials(defaultCredentials)

	return session.Must(session.NewSession(config))
}
