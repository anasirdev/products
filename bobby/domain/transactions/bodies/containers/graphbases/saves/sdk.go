package saves

import (
	"github.com/steve-care-software/products/bobby/domain/selectors"
	"github.com/steve-care-software/products/bobby/domain/transactions/bodies/containers/tables"
	"github.com/steve-care-software/products/libs/hash"
)

// NewBuilder creates a new builder instance
func NewBuilder() Builder {
	hashAdapter := hash.NewAdapter()
	return createBuilder(hashAdapter)
}

// Builder represents a transaction builder
type Builder interface {
	Create() Builder
	WithMetaData(meta tables.Transaction) Builder
	WithGraphbase(graphbase selectors.Selector) Builder
	WithParent(parent selectors.Selector) Builder
	Now() (Transaction, error)
}

// Transaction represents a create graphbase transaction
type Transaction interface {
	Hash() hash.Hash
	MetaData() tables.Transaction
	HasParent() bool
	Parent() selectors.Selector
	HasGraphbase() bool
	Graphbase() selectors.Selector
}
