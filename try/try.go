package try

func Do(f func() (interface{}, error), attempts int) (interface{}, error) {
	attemptsRemainig := attempts
	for {
		result, err := f()
		if err != nil {
			attemptsRemainig--
			if attemptsRemainig <= 0 {
				return nil, err
			}
		} else {
			return result, nil
		}
	}
}
