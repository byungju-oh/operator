package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

// DeepCopyObject returns a deep copy of the Message as a runtime.Object
func (m *Message) DeepCopyObject() runtime.Object {
	if c := m.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopy returns a deep copy of the Message
func (m *Message) DeepCopy() *Message {
	if m == nil {
		return nil
	}
	out := new(Message)
	m.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties of this object into another object of the same type
func (m *Message) DeepCopyInto(out *Message) {
	*out = *m
	out.TypeMeta = m.TypeMeta
	m.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = m.Spec
	out.Status = m.Status
}

//+kubebuilder:object:root=true

// MessageList contains a list of Message
type MessageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Message `json:"items"`
}

// DeepCopyObject returns a deep copy of the MessageList as a runtime.Object
func (ml *MessageList) DeepCopyObject() runtime.Object {
	if c := ml.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopy returns a deep copy of the MessageList
func (ml *MessageList) DeepCopy() *MessageList {
	if ml == nil {
		return nil
	}
	out := new(MessageList)
	ml.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties of this object into another object of the same type
func (ml *MessageList) DeepCopyInto(out *MessageList) {
	*out = *ml
	out.TypeMeta = ml.TypeMeta
	ml.ListMeta.DeepCopyInto(&out.ListMeta)
	if ml.Items != nil {
		in, out := &ml.Items, &out.Items
		*out = make([]Message, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func init() {
	SchemeBuilder.Register(&Message{}, &MessageList{})
}