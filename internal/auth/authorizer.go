package auth

import (
	"fmt"

	casbin "github.com/casbin/casbin/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Authorizer struct {
	enforcer *casbin.Enforcer
}

func New(model, policy string) *Authorizer {
	enforcer, err := casbin.NewEnforcer(model, policy)

	if err != nil {
		panic(err)
	}
	return &Authorizer{enforcer: enforcer}
}

func (a *Authorizer) Authorize(subject, object, action string) error {
	if ok, _ := a.enforcer.Enforce(subject, object, action); !ok {
		msg := fmt.Sprintf("%s not permitted to %s to %s", subject, action, object)
		st := status.New(codes.PermissionDenied, msg)
		return st.Err()
	}
	return nil
}
