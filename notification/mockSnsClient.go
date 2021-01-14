package notification

import (
	. "github.com/yiGmMk/pz-infra-new/logging"

	"github.com/aws/aws-sdk-go/service/sns"
)

type SnsClient interface {
	Publish(input *sns.PublishInput) (*sns.PublishOutput, error)
	CreatePlatformEndpoint(input *sns.CreatePlatformEndpointInput) (*sns.CreatePlatformEndpointOutput, error)
	DeleteEndpoint(input *sns.DeleteEndpointInput) (*sns.DeleteEndpointOutput, error)
	GetEndpointAttributes(input *sns.GetEndpointAttributesInput) (*sns.GetEndpointAttributesOutput, error)
	SetEndpointAttributes(*sns.SetEndpointAttributesInput) (*sns.SetEndpointAttributesOutput, error)
}

type MockSnsClient struct {
}

func (*MockSnsClient) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	Log.Debug("mock publish message", With("input", input))
	return &sns.PublishOutput{}, nil
}

func (*MockSnsClient) CreatePlatformEndpoint(input *sns.CreatePlatformEndpointInput) (*sns.CreatePlatformEndpointOutput, error) {
	arn := "mock_platform_arn_" + *input.Token
	Log.Debug("mock create platform endpoint.", With("input", input))
	return &sns.CreatePlatformEndpointOutput{EndpointArn: &arn}, nil
}

func (*MockSnsClient) DeleteEndpoint(input *sns.DeleteEndpointInput) (*sns.DeleteEndpointOutput, error) {
	Log.Debug("mock delete platform endpoint.", With("input", input))
	return &sns.DeleteEndpointOutput{}, nil
}

func (*MockSnsClient) GetEndpointAttributes(input *sns.GetEndpointAttributesInput) (*sns.GetEndpointAttributesOutput, error) {
	Log.Debug("mock get endpoint attributes . ", With("input", input))
	token := "mock_token"
	enabled := "true"
	return &sns.GetEndpointAttributesOutput{
		Attributes: map[string]*string{
			"Token":   &token,
			"Enabled": &enabled,
		}}, nil
}

func (*MockSnsClient) SetEndpointAttributes(input *sns.SetEndpointAttributesInput) (*sns.SetEndpointAttributesOutput, error) {
	Log.Debug("mock get endpoint attributes .", With("input", input))
	return &sns.SetEndpointAttributesOutput{}, nil
}
