package opctlutil

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/opctl/sdk-golang/model"
	nodeapi "github.com/opctl/sdk-golang/node/api/client"
	"golang.org/x/crypto/ssh"
	"log"
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
		log.Fatal(err.Error())
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

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
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}