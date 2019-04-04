/*
Copyright 2019 Zach Puckett.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package akscluster

import (
	"context"
	"fmt"
	azurev1beta1 "github.com/zachpuck/aks-engine-automation/pkg/apis/azure/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/zachpuck/aks-engine-automation/pkg/k8sutil"
	"github.com/zachpuck/aks-engine-automation/pkg/opctlutil"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AksCluster Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAksCluster{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("akscluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to AksCluster
	err = c.Watch(&source.Kind{Type: &azurev1beta1.AksCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by AksCluster - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &azurev1beta1.AksCluster{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileAksCluster{}

// ReconcileAksCluster reconciles a AksCluster object
type ReconcileAksCluster struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AksCluster object and makes changes based on the state read
// and what is in the AksCluster.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=azure.cnct.io,resources=aksclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=azure.cnct.io,resources=aksclusters/status,verbs=get;update;patch
func (r *ReconcileAksCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the AksCluster instance
	instance := &azurev1beta1.AksCluster{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	fmt.Println("status: ", instance.Status.Phase)
	if instance.Status.Phase == "" {

		// create a new opctl
		opctl := opctlutil.New("localhost")

		// generate ssh key
		sshPrivateKey, sshPublicKey, err := opctlutil.GenerateSshKey()
		if err != nil {
			return reconcile.Result{}, err
		}
		// save private key as secret
		err = k8sutil.CreateSSHSecret(r.Client, instance.Name, instance.Namespace, sshPrivateKey)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Creates a config object from AksCluster custom resource spec
		config := opctlutil.ClusterConfig{
			APIVersion: "vlabs",
			Properties: opctlutil.Properties{
				OrchestratorProfile: opctlutil.OrchestratorProfile{
					OrchestratorType:    "Kubernetes",
					OrchestratorRelease: instance.Spec.KubernetesVersion,
				},
				MasterProfile: opctlutil.MasterProfile{
					Count:     instance.Spec.MasterProfile.Count,
					DNSPrefix: instance.Spec.MasterProfile.DnsPrefix,
					VMSize:    instance.Spec.MasterProfile.VmSize,
				},
				// TODO: add for loop for multiple node pools
				AgentPoolProfiles: []opctlutil.AgentPoolProfiles{
					{
						Name:   instance.Spec.AgentPoolProfiles[0].Name,
						Count:  instance.Spec.AgentPoolProfiles[0].Count,
						VMSize: instance.Spec.AgentPoolProfiles[0].VmSize,
					},
				},
				LinuxProfile: opctlutil.LinuxProfile{
					AdminUsername: "azureuser",
					SSH: opctlutil.SSH{
						PublicKeys: []opctlutil.PublicKeys{
							{
								KeyData: sshPublicKey,
							},
						},
					},
				},
				// TODO: generate cluster specific service principal.
				// az ad sp create-for-rbac --role="Contributor" \
				// --scopes="/subscriptions/<subscriptionId>/resourceGroups/<clusterResourceGroup>"
				ServicePrincipalProfile: opctlutil.ServicePrincipalProfile{
					ClientID: instance.Spec.Credentials.LoginId,
					Secret:   instance.Spec.Credentials.LoginSecret,
				},
			},
		}

		// Starts the create-cluster operation
		opId, err := opctl.CreateCluster(opctlutil.CreateClusterInput{
			Credentials: azurev1beta1.AzureCredentials{
				TenantId:       instance.Spec.Credentials.TenantId,
				SubscriptionId: instance.Spec.Credentials.SubscriptionId,
				LoginId:        instance.Spec.Credentials.LoginId,
				LoginSecret:    instance.Spec.Credentials.LoginSecret,
			},
			Location:    instance.Spec.Location,
			ClusterName: instance.Name,
			Config:      config,
		})
		if err != nil {
			return reconcile.Result{}, err
		}

		// TODO: track opId result
		fmt.Println("opId: ", opId)

		// TODO: update status to Creating and current k8s version
		instance.Status.Phase = "Creating"
		instance.Status.K8sVersion = azurev1beta1.ClusterKubernetesVersion(instance.Spec.KubernetesVersion)
	}

	return reconcile.Result{}, nil
}
