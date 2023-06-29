package migration

// TODO: Add test cases later
func upgradeDatabaseToVersion3_3_0(fromVersion string) error {
	return updateSchemaVersion(version3_3_0)
}

func upgradeDatabaseToVersion3_3_1(fromVersion string) error {
	return updateSchemaVersion(version3_3_1)
}
