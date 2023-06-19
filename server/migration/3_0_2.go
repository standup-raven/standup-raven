package migration

func upgradeDatabaseToVersion3_0_2(fromVersion string) error {
	return updateSchemaVersion(version3_0_2)
}
