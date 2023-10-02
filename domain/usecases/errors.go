package usecases

type ErrUseCase struct {
	StatusCode int
	Err        error
}

func (e *ErrUseCase) Error() string {
	return e.Err.Error()
}

func (e *ErrUseCase) Unwrap() error {
	return e.Err
}

func (e *ErrUseCase) Is(err error) bool {
	t, ok := err.(*ErrUseCase)
	if !ok {
		return false
	}
	return t.StatusCode == e.StatusCode
}
