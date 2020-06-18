package migration

import (
	"bou.ke/monkey"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpgradeDatabaseToVersion2_0_0(t *testing.T)  {
	defer TearDown()

	updateSchemaVersionCount := 0
	monkey.Patch(updateSchemaVersion, func(version string) error {
		updateSchemaVersionCount++
		return nil
	})
	
	err := upgradeDatabaseToVersion2_0_0(version2_0_0, version3_0_0)
	assert.Nil(t, err)
	assert.Equal(t, 1, updateSchemaVersionCount)
}

func TestUpgradeDatabaseToVersion2_0_0_updateSchemaVersion_error(t *testing.T)  {
	defer TearDown()

	updateSchemaVersionCount := 0
	monkey.Patch(updateSchemaVersion, func(version string) error {
		updateSchemaVersionCount++
		return errors.New("")
	})

	err := upgradeDatabaseToVersion2_0_0(version2_0_0, version3_0_0)
	assert.NotNil(t, err)
	assert.Equal(t, 1, updateSchemaVersionCount)
}


