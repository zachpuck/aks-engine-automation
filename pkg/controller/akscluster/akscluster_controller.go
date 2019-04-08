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
	"k8s.io/apimachinery/pkg/types"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"

	opctlmodel "github.com/opctl/sdk-golang/model"
	"github.com/zachpuck/aks-engine-automation/pkg/k8sutil"
	"github.com/zachpuck/aks-engine-automation/pkg/opctlutil"
)

var log = logf.Log.WithName("controller")

const (
	ClusterPhaseNone              = ""
	ClusterPhasePending           = "Pending"
	ClusterPhaseWaitingForCluster = "Waiting"
	ClusterPhaseUpgrading         = "Upgrading"
	ClusterPhaseDeleting          = "Deleting"
	ClusterPhaseReady             = "Ready"
	ClusterPhaseFailed            = "Failed"
	ClusterPhaseError             = "Error"
	OpctlHostname                 = "localhost"
)

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
	clusterInstance := &azurev1beta1.AksCluster{}
	err := r.Get(context.TODO(), request.NamespacedName, clusterInstance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// confirm credentials secret before create cluster
	secretResults, err := k8sutil.GetSecret(r.Client, clusterInstance.Spec.Credentials, clusterInstance.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}
	clusterCredentials := azurev1beta1.AzureCredentials{
		TenantId:       string(secretResults.Data["tenantId"]),
		SubscriptionId: string(secretResults.Data["subscriptionId"]),
		LoginId:        string(secretResults.Data["loginId"]),
		LoginSecret:    string(secretResults.Data["loginSecret"]),
	}

	// Create cluster
	if clusterInstance.Status.Phase == "" {
		log.Info("creating", "Cluster:", clusterInstance.Name)

		// create a new opctl
		opctl := opctlutil.New("localhost")

		// generate ssh key
		sshPrivateKey, sshPublicKey, err := opctlutil.GenerateSshKey()
		if err != nil {
			return reconcile.Result{}, err
		}
		// save private key as secret
		err = k8sutil.CreateSSHSecret(r.Client, clusterInstance.Name, clusterInstance.Namespace, sshPrivateKey)
		if err != nil {
			return reconcile.Result{}, err
		}

		// multiple agent node pools
		var newAgentPoolProfiles []opctlutil.AgentPoolProfiles
		for i := range clusterInstance.Spec.AgentPoolProfiles {
			newAgentPoolProfiles = append(newAgentPoolProfiles, opctlutil.AgentPoolProfiles{
				Name:   clusterInstance.Spec.AgentPoolProfiles[i].Name,
				Count:  clusterInstance.Spec.AgentPoolProfiles[i].Count,
				VMSize: clusterInstance.Spec.AgentPoolProfiles[i].VmSize,
			})
		}

		// Creates a config object from AksCluster custom resource spec
		config := opctlutil.ClusterConfig{
			APIVersion: "vlabs",
			Properties: opctlutil.Properties{
				OrchestratorProfile: opctlutil.OrchestratorProfile{
					OrchestratorType:    "Kubernetes",
					OrchestratorRelease: clusterInstance.Spec.KubernetesVersion,
				},
				MasterProfile: opctlutil.MasterProfile{
					Count:     clusterInstance.Spec.MasterProfile.Count,
					DNSPrefix: clusterInstance.Spec.MasterProfile.DnsPrefix,
					VMSize:    clusterInstance.Spec.MasterProfile.VmSize,
				},
				AgentPoolProfiles: newAgentPoolProfiles,
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
					ClientID: clusterCredentials.LoginId,
					Secret:   clusterCredentials.LoginSecret,
				},
			},
		}

		// Starts the create-cluster operation
		createClusterResult, err := opctl.CreateCluster(opctlutil.CreateClusterInput{
			Credentials: clusterCredentials,
			Location:    clusterInstance.Spec.Location,
			ClusterName: clusterInstance.Name,
			Config:      config,
		})
		if err != nil {
			return reconcile.Result{}, err
		}

		matchPattern := regexp.MustCompile("^[a-zA-Z0-9]*$")
		if !matchPattern.MatchString(createClusterResult.OpId) {

			statusUpdate := ClusterInstanceUpdates{
				Name:      clusterInstance.Name,
				Namespace: clusterInstance.Namespace,
				Phase:     ClusterPhaseError,
			}
			err = r.updateClusterInstance(statusUpdate)
			if err != nil {
				return reconcile.Result{}, fmt.Errorf("failed to update cluster status: %v", err)
			}

			return reconcile.Result{}, fmt.Errorf(
				"failed to start operation for cluster -->%v<-- with message: %v",
				clusterInstance.Name, createClusterResult.OpId,
			)
		}

		createClusterUpdates := ClusterInstanceUpdates{
			Name:            clusterInstance.Name,
			Namespace:       clusterInstance.Namespace,
			Phase:           ClusterPhasePending,
			K8sVersion:      clusterInstance.Spec.KubernetesVersion,
			AnnotationKey:   "createClusterOpId",
			AnnotationValue: createClusterResult.OpId,
		}
		err = r.updateClusterInstance(createClusterUpdates)
		if err != nil {
			return reconcile.Result{}, err
		}

		// track operation result
		go r.ResolveOperation(ResolveOperationInput{
			Name:      clusterInstance.Name,
			Namespace: clusterInstance.Namespace,
			OpId:      createClusterResult.OpId,
			StartTime: time.Now().UTC(),
		})
	}

	return reconcile.Result{}, nil
}

type ClusterInstanceUpdates struct {
	Name            string
	Namespace       string
	Phase           string
	K8sVersion      string
	AnnotationKey   string
	AnnotationValue string
}

// updateClusterInstance update status fields of the aks cluster instance object and emits events
func (r *ReconcileAksCluster) updateClusterInstance(input ClusterInstanceUpdates) error {
	// Fetch the AksCluster instance
	clusterInstance := &azurev1beta1.AksCluster{}

	err := r.Get(
		context.Background(),
		client.ObjectKey(
			types.NamespacedName{
				Name:      input.Name,
				Namespace: input.Namespace,
			}),
		clusterInstance,
	)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Error(err, "Object not found", "cluster", clusterInstance.Name)
		}
	}

	clusterInstanceCopy := clusterInstance.DeepCopy()

	// update status
	clusterInstanceCopy.Status.Phase = azurev1beta1.ClusterStatusPhase(input.Phase)
	// update k8s version
	if input.K8sVersion != "" {
		clusterInstanceCopy.Status.KubernetesVersion = azurev1beta1.ClusterKubernetesVersion(input.K8sVersion)
	}
	// update annotation
	if input.AnnotationKey != "" {
		clusterInstanceCopy.Annotations[input.AnnotationKey] = input.AnnotationValue
	}

	// update status of AksCluster resource
	err = r.Client.Update(context.Background(), clusterInstanceCopy)
	if err != nil {
		return fmt.Errorf(
			"failed to update aks cluster resource status fields for cluster %v: %v",
			clusterInstanceCopy.Name, err,
		)
	}
	return nil
}

type ResolveOperationInput struct {
	Name      string
	Namespace string
	OpId      string
	StartTime time.Time
}

// ResolveOperation waits for the provided operation to complete and update the aks cluster resource status phase
func (r *ReconcileAksCluster) ResolveOperation(input ResolveOperationInput) {
	// create a new opctl
	opctl := opctlutil.New(OpctlHostname)

	result, err := opctl.GetOpEvents(opctlutil.GetOpEventsInput{
		OpId:      input.OpId,
		StartTime: input.StartTime,
	})
	if err != nil {
		log.Error(err, "failed to resolve operation for", "cluster", input.Name)
	}

	log.Info("Operation complete with", "outcome", result.Outcome)

	// update cluster status
	update := ClusterInstanceUpdates{
		Name:      input.Name,
		Namespace: input.Namespace,
	}
	if result.Outcome == opctlmodel.OpOutcomeSucceeded {
		update.Phase = ClusterPhaseReady
	} else {
		update.Phase = ClusterPhaseFailed
	}
	err = r.updateClusterInstance(update)
	if err != nil {
		log.Error(err, "failed to update status phase on", "cluster", input.Name)
	}
}
