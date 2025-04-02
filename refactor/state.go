package compare

// State represents the state of a running comparison.
type State struct {
	Path Path
}

// Difference creates a new [Difference] for the given arguments.
func (s *State) Difference(msg string, v1, v2 any) Difference {
	return Difference{
		Path:    s.Path.Clone(),
		Message: msg,
		V1:      v1,
		V2:      v2,
	}
}

// Result creates a new [Result] with a single [Difference] for the given arguments.
func (s *State) Result(msg string, v1, v2 any) (Result, bool) {
	return Result{s.Difference(msg, v1, v2)}, true
}
