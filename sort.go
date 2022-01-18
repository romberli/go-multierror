package multierror

// Len implements sort.Interface function for length
func (err Error) Len() int {
	return len(err.Errs)
}

// Swap implements sort.Interface function for swapping elements
func (err Error) Swap(i, j int) {
	err.Errs[i], err.Errs[j] = err.Errs[j], err.Errs[i]
}

// Less implements sort.Interface function for determining order
func (err Error) Less(i, j int) bool {
	return err.Errs[i].Error() < err.Errs[j].Error()
}
