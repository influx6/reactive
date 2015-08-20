package flow

import (
	"fmt"
	"reflect"
	"time"

	"github.com/influx6/flux"
)

type (
	//Immutable defines an interface method rules for immutables types. All types meeting this rule must be single type values
	Immutable interface {
		Value() interface{}
		Clone() Immutable
		set(interface{}) bool
	}

	//Mutation defines the basic operation change that occurs with an object
	Mutation struct {
		Immutable
		timestamp time.Time
	}

	//ReactiveObserver defines a basic reactive value
	ReactiveObserver struct {
		flux.Stacks
		data *Mutation
	}
)

//NewMutation returns a new mutation marked with a stamp
func NewMutation(m Immutable) *Mutation {
	return &Mutation{
		Immutable: m,
		timestamp: time.Now(),
	}
}

const (
	//ErrUnacceptedTypeMessage defines the message for types that are not part of the basic units/types in go
	ErrUnacceptedTypeMessage = "Type %s is not acceptable"
)

//MakeType validates accepted types and returns the (Immutable, error)
func MakeType(val interface{}) (Immutable, error) {

	switch reflect.TypeOf(val).Kind() {
	case reflect.Struct:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "struct")
	case reflect.Map:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "map")
	case reflect.Array:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "array")
	case reflect.Slice:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "slice")
	}

	return nil, nil
}

//Reactive returns a new Reactive instance
func Reactive(m interface{}) (*ReactiveObserver, error) {
	var im Immutable
	var err error

	if im, err = MakeType(m); err != nil {
		return nil, err
	}

	return ReactiveFrom(im), nil
}

//ReactiveFrom returns a new Reactive instance
func ReactiveFrom(m Immutable) *ReactiveObserver {
	return &ReactiveObserver{
		Stacks: flux.IdentityStack(),
		data:   NewMutation(m),
	}
}

//Set resets the value of the object
func (r *ReactiveObserver) Set(ndata interface{}) {
	clone := r.data.Clone()

	//can we make the change or his this change proper
	if !clone.set(ndata) {
		return
	}

	mut := NewMutation(clone)
	r.Call(mut)
	r.data = mut
}

//Get returns the internal value
func (r *ReactiveObserver) Get() interface{} {
	return r.data.Value()
}
