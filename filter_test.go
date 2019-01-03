package gopiper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostAdd(t *testing.T) {
	p := &PipeItem{}
	r, e := p.CallFilter(`i`, `postadd(s)`)
	if e != nil {
		t.Fatal(e)
	}
	assert.Equal(t, `is`, r)
}
