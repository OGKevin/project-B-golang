package acl

import (
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/persist"
)

func NewEnforcer(a persist.Adapter) *casbin.Enforcer {
	// Or you can use an existing DB "abc" like this:
	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	// a := xormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/abc", true)

	 return casbin.NewEnforcer(casbin.NewModel(conf), a)
	//
	//// Load the policy from DB.
	//e.LoadPolicy()
	//
	//// Check the permission.
	//e.Enforce("alice", "data1", "read")
	//e.AddPolicy()
	//
	//// Modify the policy.
	//// e.AddPolicy(...)
	//// e.RemovePolicy(...)
	//
	//// Save the policy back to DB.
	//e.SavePolicy()
}
