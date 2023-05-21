package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashKey(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a := &String{Value: "1"}
		b := &String{Value: "1"}
		assert.Equal(t, a.HashKey(), b.HashKey())
		c := &String{Value: "2"}
		d := &String{Value: "2"}
		assert.Equal(t, c.HashKey(), d.HashKey())
		assert.NotEqual(t, a.HashKey(), c.HashKey())
	})
	t.Run("int", func(t *testing.T) {
		a := &Integer{Value: 1}
		b := &Integer{Value: 1}
		assert.Equal(t, a.HashKey(), b.HashKey())
		c := &Integer{Value: 2}
		d := &Integer{Value: 2}
		assert.Equal(t, c.HashKey(), d.HashKey())
		assert.NotEqual(t, a.HashKey(), c.HashKey())
	})
	t.Run("bool", func(t *testing.T) {
		a := &Boolean{Value: false}
		b := &Boolean{Value: false}
		assert.Equal(t, a.HashKey(), b.HashKey())
		c := &Boolean{Value: true}
		d := &Boolean{Value: true}
		assert.Equal(t, c.HashKey(), d.HashKey())
		assert.NotEqual(t, a.HashKey(), c.HashKey())
	})
}
