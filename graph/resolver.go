package graph

import (
	"github.com/gobuffalo/pop/v6"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	db *pop.Connection
}
