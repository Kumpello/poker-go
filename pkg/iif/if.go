package iif

func IfElse[T any](cond bool, then T, elze T) T { // nolint:ireturn // nolint:nolintlint
	if cond {
		return then
	} else { // nolint:revive // that's fine
		return elze
	}
}
