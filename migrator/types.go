package migrator

type Migration interface {
	ID() int64
}
