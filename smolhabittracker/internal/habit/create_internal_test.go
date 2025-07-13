package habit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_validateAndFillDetails(t *testing.T) {
	t.Parallel()

	t.Run("Full", testValidateAndFillDetailsFull)
	t.Run("Partial", testValidateAndFillDetailsPartial)
	t.Run("SpaceName", testValidateAndFillDetailsSpaceName)
}

func testValidateAndFillDetailsFull(t *testing.T) {
	t.Parallel()

	h := Habit{ID: "1", Name: "swim", WeeklyFrequency: 3, CreationTime: time.Now()}

	got, err := validateAndFillDetails(h)
	require.NoError(t, err)
	assert.Equal(t, h, got)
}

func testValidateAndFillDetailsPartial(t *testing.T) {
	t.Parallel()

	h := Habit{Name: "run", WeeklyFrequency: 6}

	got, err := validateAndFillDetails(h)
	require.NoError(t, err)
	assert.Equal(t, h.Name, got.Name)
	assert.Equal(t, h.WeeklyFrequency, got.WeeklyFrequency)
	assert.NotEmpty(t, got.ID)
	assert.NotEmpty(t, got.CreationTime)
}

func testValidateAndFillDetailsSpaceName(t *testing.T) {
	t.Parallel()

	h := Habit{Name: "    "}

	_, err := validateAndFillDetails(h)
	assert.ErrorAs(t, err, &InvalidInputError{})
}
