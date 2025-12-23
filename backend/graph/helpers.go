package graph

// ptrToInt безопасно разыменовывает указатель на int.
func ptrToInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}
