package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Nescaffe",
		Price: 1.00,
		SKU:   "abs-abc-def",
	}

	val := NewValidation()
	err := val.Validate(p)

	assert.Len(t, err, 1)
}
