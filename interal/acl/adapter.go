package acl

import "github.com/casbin/casbin/model"

type adapter struct {

}

func (*adapter) LoadPolicy(model model.Model) error {
	panic("implement me")
}

func (*adapter) SavePolicy(model model.Model) error {
	panic("implement me")
}

func (*adapter) AddPolicy(sec string, ptype string, rule []string) error {
	panic("implement me")
}

func (*adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	panic("implement me")
}

func (*adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	panic("implement me")
}

