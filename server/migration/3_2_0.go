package migration

func upgradeDatabaseToVersion3_2_0(fromVersion string) error {
	if UpdateErr := updateSchemaVersion(version3_2_0); UpdateErr != nil {
		return UpdateErr
	}
	return nil
}
