package lib

func MergeValues[K comparable, V any](m1 map[K][]V, m2 map[K][]V) {
	for k, v := range m2 {
		if _, ok := m1[k]; !ok {
			m1[k] = []V{}
		}
		m1[k] = append(m1[k], v...)
	}
}
