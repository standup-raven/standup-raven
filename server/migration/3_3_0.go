package migration

func upgradeDatabaseToVersion3_3_0(fromVersion string) error {
	return updateSchemaVersion(version3_3_0)
}
