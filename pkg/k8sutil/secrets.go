package k8sutil

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const PrivateKeySuffix = "privatekey"

func GetSecret(c client.Client, name string, namespace string) (corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := c.Get(
		context.Background(),
		client.ObjectKey{Namespace: namespace, Name: name},
		secret)
	if err != nil {
		return *secret, err
	}

	return *secret, err
}

func CreateSecret(c client.Client, secret *corev1.Secret) error {
	err := c.Create(
		context.Background(),
		secret)
	if err != nil {
		glog.Errorf("failed to create Secret %s: %q", secret.Name, err)
		return err
	}

	return nil
}

func CreateSSHSecret(c client.Client, ownerRef []v1.OwnerReference, clusterName string, namespace string, privateKey []byte) error {
	secretName := fmt.Sprintf("%s-%s", clusterName, PrivateKeySuffix)
	dataMap := make(map[string][]byte)
	dataMap[corev1.SSHAuthPrivateKey] = privateKey

	secret, err := GetSecret(c, secretName, namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			newSecret := &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Name:      secretName,
					Namespace: namespace,
				},
				Type: corev1.SecretTypeSSHAuth,
				Data: dataMap,
			}

			newSecret.OwnerReferences = ownerRef
			err = CreateSecret(c, newSecret)
			if err != nil {
				glog.Errorf("failed to create ssh secret %s: %q", secret.GetName(), err)
				return err
			}
		} else {
			glog.Errorf("failed to query for secret %s: %q", secretName, err)
		}
	} else {
		if !reflect.DeepEqual(secret.Data, dataMap) {
			secret.Data = dataMap
			err = c.Update(context.Background(), &secret)
			if err != nil {
				glog.Errorf("failed to update secret %s: %q", secretName, err)
			}
		}
	}
	return nil
}
