package migrator

type Entity interface {
	ID() int64
	CompareTo(dst Entity) bool
}
