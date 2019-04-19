package opctlutil

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/opctl/sdk-golang/model"
	nodeapi "github.com/opctl/sdk-golang/node/api/client"
	"golang.org/x/crypto/ssh"
	"net/url"
	"os"
)

type OPCTL struct {
	Client nodeapi.Client
}

// New creates a new opctl client
func New(hostName string) *OPCTL {
	return &OPCTL{
		Client: nodeapi.New(url.URL{
			Scheme: "http",
			Host:   hostName + ":42224",
			Path:   "api",
		},
			nil,
		),
	}
}

func (o *OPCTL) CreateCluster(input CreateClusterInput) (CreateClusterOutput, error) {
	// create config object
	var configObject map[string]interface{}

	encodedConfig, err := json.Marshal(input.Config)
	if err != nil {
		return CreateClusterOutput{}, fmt.Errorf("failed to encoded config: %v", err)
	}
	if err := json.Unmarshal(encodedConfig, &configObject); err != nil {
		return CreateClusterOutput{}, fmt.Errorf("failed to create configObject: %v", err)
	}

	opId, err := o.Client.StartOp(
		context.Background(),
		model.StartOpReq{
			Args: map[string]*model.Value{
				"subscriptionId": {String: to.StringPtr(input.Credentials.SubscriptionId)},
				"loginId":        {String: to.StringPtr(input.Credentials.LoginId)},
				"loginSecret":    {String: to.StringPtr(input.Credentials.LoginSecret)},
				"loginTenantId":  {String: to.StringPtr(input.Credentials.TenantId)},
				"location":       {String: to.StringPtr(input.Location)},
				"clusterName":    {String: to.StringPtr(input.ClusterName)},
				"config":         {Object: configObject},
				// Provided as part of the deployment
				"storageAccountName":              {String: to.StringPtr(os.Getenv("AKS_ENGINE_STORAGE_ACCOUNT_NAME"))},
				"storageAccountResourceGroupName": {String: to.StringPtr(os.Getenv("AKS_ENGINE_STORAGE_ACCOUNT_GROUP"))},
			},
			Op: model.StartOpReqOp{
				Ref: os.Getenv("OPERATIONS_PKG_PATH") + "/create-cluster",
			},
		})
	if err != nil {
		return CreateClusterOutput{}, fmt.Errorf("failed to start create-cluster op: %v", err)
	}
	return CreateClusterOutput{
		OpId: opId,
	}, nil
}

func GenerateSshKey() ([]byte, string, error) {
	bitSize := 2048

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		return nil, "", err
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, "", err
	}

	privateKeyBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	return privateKeyBytes, string(publicKeyBytes), nil
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(publickey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(publickey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}

// GetPublicKey from private key
func GetPublicKeyFromPrivate(privateKey string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))

	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	publicKeyBytes, _ := generatePublicKey(&key.PublicKey)

	return string(publicKeyBytes), nil
}

// GetOpEvents starts event loop while operation is running until completion
func (o *OPCTL) GetOpEvents(input GetOpEventsInput) (GetOpEventsOutput, error) {
	eventChannel, err := o.Client.GetEventStream(
		&model.GetEventStreamReq{
			Filter: model.EventFilter{
				Roots: []string{input.OpId},
				Since: &input.StartTime,
			},
		})
	if err != nil {
		return GetOpEventsOutput{}, err
	}

	for {
		select {

		case event, isEventChannelOpen := <-eventChannel:
			if !isEventChannelOpen {
				return GetOpEventsOutput{}, err
			}

			if event.OpEnded != nil {
				if event.OpEnded.OpID == input.OpId {
					switch event.OpEnded.Outcome {
					case model.OpOutcomeSucceeded:
						return GetOpEventsOutput{Outcome: model.OpOutcomeSucceeded}, nil
					case model.OpOutcomeKilled:
						return GetOpEventsOutput{Outcome: model.OpOutcomeKilled}, nil
					default:
						return GetOpEventsOutput{Outcome: model.OpOutcomeFailed}, nil
					}
				}
			}
		}
	}
}

// DeleteCluster deletes an aks-engine cluster
func (o *OPCTL) DeleteCluster(input DeleteClusterInput) (DeleteClusterOutput, error) {
	opId, err := o.Client.StartOp(
		context.Background(),
		model.StartOpReq{
			Args: map[string]*model.Value{
				"subscriptionId":                  {String: to.StringPtr(input.Credentials.SubscriptionId)},
				"loginId":                         {String: to.StringPtr(input.Credentials.LoginId)},
				"loginSecret":                     {String: to.StringPtr(input.Credentials.LoginSecret)},
				"loginTenantId":                   {String: to.StringPtr(input.Credentials.TenantId)},
				"location":                        {String: to.StringPtr(input.Location)},
				"clusterName":                     {String: to.StringPtr(input.ClusterName)},
				"storageAccountName":              {String: to.StringPtr(os.Getenv("AKS_ENGINE_STORAGE_ACCOUNT_NAME"))},
				"storageAccountResourceGroupName": {String: to.StringPtr(os.Getenv("AKS_ENGINE_STORAGE_ACCOUNT_GROUP"))},
			},
			Op: model.StartOpReqOp{
				Ref: os.Getenv("OPERATIONS_PKG_PATH") + "/delete-cluster",
			},
		})
	if err != nil {
		return DeleteClusterOutput{}, fmt.Errorf("failed to start delete-cluster op: %v", err)
	}

	return DeleteClusterOutput{
		OpId: opId,
	}, nil
}

// AddNodePoolGroup adds a node pool group to the cluster
func (o *OPCTL) AddNodePoolGroup(input AddNodePoolGroupInput) (AddNodePoolGroupOutput, error) {
	// create config object
	var configObject map[string]interface{}

	encodedConfig, err := json.Marshal(input.Config)
	if err != nil {
		return AddNodePoolGroupOutput{}, fmt.Errorf("failed to encoded config: %v", err)
	}
	if err := json.Unmarshal(encodedConfig, &configObject); err != nil {
		return AddNodePoolGroupOutput{}, fmt.Errorf("failed to create configObject: %v", err)
	}

	opId, err := o.Client.StartOp(
		context.Background(),
		model.StartOpReq{
			Args: map[string]*model.Value{
				"subscriptionId": {String: to.StringPtr(input.Credentials.SubscriptionId)},
				"loginId":        {String: to.StringPtr(input.Credentials.LoginId)},
				"loginSecret":    {String: to.StringPtr(input.Credentials.LoginSecret)},
				"loginTenantId":  {String: to.StringPtr(input.Credentials.TenantId)},
				"clusterName":    {String: to.StringPtr(input.ClusterName)},
				"config":         {Object: configObject},
				// Provided as part of the deployment
				"storageAccountName":              {String: to.StringPtr(os.Getenv("AKS_ENGINE_STORAGE_ACCOUNT_NAME"))},
				"storageAccountResourceGroupName": {String: to.StringPtr(os.Getenv("AKS_ENGINE_STORAGE_ACCOUNT_GROUP"))},
			},
			Op: model.StartOpReqOp{
				Ref: os.Getenv("OPERATIONS_PKG_PATH") + "/add-node-pool",
			},
		})
	if err != nil {
		return AddNodePoolGroupOutput{}, fmt.Errorf("failed to start add-node-pool op: %v", err)
	}

	return AddNodePoolGroupOutput{
		OpId: opId,
	}, nil
}

// RemoveNodePoolGroup removes a node pool group from the cluster
func (o *OPCTL) RemoveNodePoolGroup(input RemoveNodePoolGroupInput) (RemoveNodePoolGroupOutput, error) {
	// create config object
	var configObject map[string]interface{}

	encodedConfig, err := json.Marshal(input.Config)
	if err != nil {
		return RemoveNodePoolGroupOutput{}, fmt.Errorf("failed to encoded config: %v", err)
	}
	if err := json.Unmarshal(encodedConfig, &configObject); err != nil {
		return RemoveNodePoolGroupOutput{}, fmt.Errorf("failed to create configObject: %v", err)
	}

	opId, err := o.Client.StartOp(
		context.Background(),
		model.StartOpReq{
			Args: map[string]*model.Value{
				"subscriptionId": {String: to.StringPtr(input.Credentials.SubscriptionId)},
				"loginId":        {String: to.StringPtr(input.Credentials.LoginId)},
				"loginSecret":    {String: to.StringPtr(input.Credentials.LoginSecret)},
				"loginTenantId":  {String: to.StringPtr(input.Credentials.TenantId)},
				"location":       {String: to.StringPtr(input.Location)},
				"clusterName":    {String: to.StringPtr(input.ClusterName)},
				"config":         {Object: configObject},
				// Provided as part of the deployment
				"storageAccountName":              {String: to.StringPtr(os.Getenv("AKS_ENGINE_STORAGE_ACCOUNT_NAME"))},
				"storageAccountResourceGroupName": {String: to.StringPtr(os.Getenv("AKS_ENGINE_STORAGE_ACCOUNT_GROUP"))},
			},
			Op: model.StartOpReqOp{
				Ref: os.Getenv("OPERATIONS_PKG_PATH") + "/remove-node-pool",
			},
		})
	if err != nil {
		return RemoveNodePoolGroupOutput{}, fmt.Errorf("failed to start remove-node-pool op: %v", err)
	}

	return RemoveNodePoolGroupOutput{
		OpId: opId,
	}, nil
}
