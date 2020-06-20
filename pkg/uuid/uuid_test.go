package uuid_test

import (
	"testing"

	"github.com/xabi93/lana-test/pkg/uuid"

	gUuid "github.com/google/uuid"

	"github.com/stretchr/testify/require"
)

func TestNewGeneratesNewValidUUID(t *testing.T) {
	require.NotPanics(t, func() { gUuid.MustParse(uuid.New()) })
}
