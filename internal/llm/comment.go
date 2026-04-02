package cloud_agent

func Async[T any](f func() (T, error)) (chan T, chan error) {
	resChan := make(chan T, 1)
	errChan := make(chan error, 1)

	go func() {
		res, err := f()
		if err != nil {
			errChan <- err
			return
		}
		resChan <- res
	}()

	return resChan, errChan
}
