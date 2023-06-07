package flatr

// Option change default Flattener behaviour
type Option func(f *Flattener)

// WithPrefix add a prefix to flatted keys
func WithPrefix(prefix string) Option {
	return func(f *Flattener) {
		f.prefix = prefix
	}
}

// WithSeparator change default separator of flatted keys
func WithSeparator(separator string) Option {
	return func(f *Flattener) {
		f.separator = separator
	}
}

// AddTransformer add a transformer called to all entries in the flattening process
func AddTransformer(tr Transformer) Option {
	return func(f *Flattener) {
		f.transformers = append(f.transformers, tr)
	}
}

// AddScopedTransformer add a transfomer called only in nodes that match the specific key
func AddScopedTransformer(key string, tr Transformer) Option {
	return func(f *Flattener) {
		f.scopedTrasformers[key] = tr
	}
}
