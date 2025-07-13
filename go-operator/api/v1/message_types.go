package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MessageSpec defines the desired state of Message
type MessageSpec struct {
	// Text is the message content to be processed
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Text string `json:"text"`
}

// MessageStatus defines the observed state of Message
type MessageStatus struct {
	// Phase represents the current phase of the Message
	Phase string `json:"phase,omitempty"`
	
	// LastUpdated is the timestamp when the message was last processed
	LastUpdated metav1.Time `json:"lastUpdated,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Text",type="string",JSONPath=".spec.text"
//+kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Message is the Schema for the messages API
type Message struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MessageSpec   `json:"spec,omitempty"`
	Status MessageStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MessageList contains a list of Message
type MessageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Message `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Message{}, &MessageList{})
}