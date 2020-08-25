package migration

func upgradeDatabaseToVersion3_0_1(fromVersion string) error {
	if UpdateErr := updateSchemaVersion(version3_0_1); UpdateErr != nil {
		return UpdateErr
	}
	return nil
}
