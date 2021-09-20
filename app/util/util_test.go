package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatPort(t *testing.T) {
	require.Equal(t, ":3000", FormatPort(3000))
	require.Equal(t, ":6739", FormatPort(6739))
	require.Equal(t, ":8080", FormatPort(8080))
}
