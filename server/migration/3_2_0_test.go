package migration

import (
	"errors"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestUpgradeDatabaseToVersion3_2_0(t *testing.T) {
	defer TearDown()

	updateSchemaVersionCount := 0
	monkey.Patch(updateSchemaVersion, func(version string) error {
		updateSchemaVersionCount++
		return nil
	})

	err := upgradeDatabaseToVersion3_2_0(version3_1_0)
	assert.Nil(t, err)
	assert.Equal(t, 1, updateSchemaVersionCount)

	err = upgradeDatabaseToVersion3_2_0(version3_1_1)
	assert.Nil(t, err)
	assert.Equal(t, 2, updateSchemaVersionCount)

	monkey.Patch(updateSchemaVersion, func(version string) error {
		updateSchemaVersionCount++
		return errors.New("simulated error")
	})

	err = upgradeDatabaseToVersion3_2_0(version3_1_1)
	assert.NotNil(t, err)
	assert.Equal(t, 3, updateSchemaVersionCount)
}

func TestUpgradeDatabaseToVersion3_2_1(t *testing.T) {
	defer TearDown()

	updateSchemaVersionCount := 0
	monkey.Patch(updateSchemaVersion, func(version string) error {
		updateSchemaVersionCount++
		return nil
	})

	err := upgradeDatabaseToVersion3_2_0(version3_1_0)
	assert.Nil(t, err)
	assert.Equal(t, 1, updateSchemaVersionCount)

	err = upgradeDatabaseToVersion3_2_0(version3_1_1)
	assert.Nil(t, err)
	assert.Equal(t, 2, updateSchemaVersionCount)

	monkey.Patch(updateSchemaVersion, func(version string) error {
		updateSchemaVersionCount++
		return errors.New("simulated error")
	})

	err = upgradeDatabaseToVersion3_2_0(version3_1_1)
	assert.NotNil(t, err)
	assert.Equal(t, 3, updateSchemaVersionCount)
}
