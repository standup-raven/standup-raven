package migration

func upgradeDatabaseToVersion3_0_1(fromVersion string) error {
	return updateSchemaVersion(version3_0_1)
}
