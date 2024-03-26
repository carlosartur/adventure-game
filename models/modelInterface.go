package models

type ModelInterface interface {
	Create() ModelInterface
	Update() ModelInterface
	Retrieve() []ModelInterface
	Delete()
	GetTableName() string
}
