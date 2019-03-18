package acl

import (
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/persist"
)

func NewEnforcer(a persist.Adapter) *casbin.Enforcer {
	 return casbin.NewEnforcer(casbin.NewModel(conf), a)
}
