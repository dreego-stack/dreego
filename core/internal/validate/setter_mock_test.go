package validate

type testSetter struct {
	data map[string]any
}

func (s *testSetter) Set(key string, value any) {
	if s.data == nil {
		s.data = make(map[string]any)
	}
	s.data[key] = value
}
