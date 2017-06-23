package compare

var (
	// MaxSliceDifferences is the maximum number of differences for a slice.
	// If the value is reached, the comparison is stopped for the current slice.
	// It is also used for array.
	// Set to 0 disables it.
	MaxSliceDifferences = 10
)
