package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	myorgv1 "github.com/yourname/go-operator/api/v1"
)

// MessageReconciler reconciles a Message object
type MessageReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=myorg.dev,resources=messages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myorg.dev,resources=messages/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myorg.dev,resources=messages/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *MessageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Message instance
	var message myorgv1.Message
	if err := r.Get(ctx, req.NamespacedName, &message); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("🗑️ Message resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "❌ Failed to get Message")
		return ctrl.Result{}, err
	}

	logger.Info("🔄 Reconciling Message", 
		"name", message.Name, 
		"namespace", message.Namespace, 
		"text", message.Spec.Text)

	// Update status
	message.Status.Phase = "Processed"
	message.Status.LastUpdated = metav1.Now()

	if err := r.Status().Update(ctx, &message); err != nil {
		logger.Error(err, "❌ Failed to update Message status")
		return ctrl.Result{}, err
	}

	// Record an event
	r.Recorder.Event(&message, "Normal", "Processed", 
		"Message processed successfully: "+message.Spec.Text)

	logger.Info("✅ Message processed successfully", 
		"text", message.Spec.Text)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MessageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myorgv1.Message{}).
		Complete(r)
}