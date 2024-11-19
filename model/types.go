package model

type Registry interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...ModelOption) (*Model, error)
}

type ModelOption func(m *Model) error

type TableName interface {
	TableName() string
}
