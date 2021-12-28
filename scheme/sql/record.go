package sql


type SqlConverter interface {
	Convert(scheme []Field, data map[string]string) (SqlConverter, error)
}