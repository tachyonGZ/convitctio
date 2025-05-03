package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitPath(t *testing.T) {
	asserts := assert.New(t)

	asserts.Equal([]string{"/"}, SplitPath("/"))
	asserts.Equal([]string{"/", "123/", "321/"}, SplitPath("/123/321"))
	asserts.Equal([]string{"/", "123/", "321/"}, SplitPath("/123/321/"))
}
