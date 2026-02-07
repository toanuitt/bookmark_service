package common

// HandleError panics if the given error is not nil.
//
// This function is intended for situations where an error is considered
// unrecoverable and the application should fail fast (e.g., during startup
// or critical initialization). If err is nil, HandleError does nothing.
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
