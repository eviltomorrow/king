package sqlutil

type condition interface {
	SQL() (string, []interface{})
	SetPrefix(Type) condition
}
