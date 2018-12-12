package replacer

// Replacer provides generic interface for replacing source
// based on patterns
type Replacer interface {
	Replace(source interface{}) (interface{}, bool)
}
