package migration

func upgradeDatabaseToVersion2_0_0(fromVersion, toVersion string) error {
	if UpdateErr := updateSchemaVersion(version2_0_0); UpdateErr != nil {
		return UpdateErr
	}
	return nil
}
