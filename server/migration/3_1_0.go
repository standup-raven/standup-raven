package migration

func upgradeDatabaseToVersion3_1_0(fromVersion string) error {
	return updateSchemaVersion(version3_1_0)
}

func upgradeDatabaseToVersion3_1_1(fromVersion string) error {
	return updateSchemaVersion(version3_1_1)
}
