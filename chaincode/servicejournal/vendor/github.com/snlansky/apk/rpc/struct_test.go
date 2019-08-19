package rpc

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestNewInternalError(t *testing.T) {
	err := errors.New("oops!")

	ierr :=NewInternalError(err, "unknown error")

	assert.Equal(t, ierr.Error(), strings.Join([]string{"unknown error", "oops!"}, ","))
	assert.Error(t, ierr)
	assert.Equal(t, ierr.External(), "unknown error")
}
