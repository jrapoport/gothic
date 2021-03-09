package migration

import "gorm.io/gorm"

// ModelA is a model for tests.
type ModelA struct {
	gorm.Model
	Value string
}

// ModelB is a model with indexes for tests.
type ModelB struct {
	gorm.Model
	Value string `gorm:"index:idx_value"`
}

// ModelBIndex is the name of the index for ModelB
const ModelBIndex = "idx_value"
