package storage

var migrations []interface{}

func AddMigration(obj interface{}) {
	migrations = append(migrations, obj)
}
