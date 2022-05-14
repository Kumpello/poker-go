package iif

func IfElse[T any](cond bool, then T, elze T) T { // nolint:ireturn // nolint:nolintlint
	if cond {
		return then
	} else { // nolint:revive // that's fine
		return elze
	}
}

func IfNil[T any](what *T, ifNil T) T {
	if what == nil {
		return ifNil
	}

	return *what
}

func EmptyIfNil[T any](what *T) T {
	if what == nil {
		var t T
		return t
	}

	return *what
}
