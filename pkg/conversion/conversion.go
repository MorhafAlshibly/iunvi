package conversion

func PointerToValue[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

func ValueToPointer[T any](value T) *T {
	return &value
}
