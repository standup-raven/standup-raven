package migration

func upgradeDatabaseToVersion3_0_2(fromVersion string) error {
	if UpdateErr := updateSchemaVersion(version3_0_2); UpdateErr != nil {
		return UpdateErr
	}
	return nil
}
