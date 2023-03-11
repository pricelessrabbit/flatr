package flatr

// Option change default Flatter behaviour
type Option func(f *Flatter)

// WithPrefix add a prefix to flatted keys
func WithPrefix(prefix string) Option {
	return func(f *Flatter) {
		f.prefix = prefix
	}
}

// WithSeparator change default separator of flatted keys
func WithSeparator(separator string) Option {
	return func(f *Flatter) {
		f.separator = separator
	}
}

// AddTransformer add a transformer called to all entries in the flattening process
func AddTransformer(tr Transformer) Option {
	return func(f *Flatter) {
		f.transformers = append(f.transformers, tr)
	}
}

// AddScopedTransformer add a transfomer called only in nodes that match the specific key
func AddScopedTransformer(key string, tr Transformer) Option {
	return func(f *Flatter) {
		f.scopedTrasformers[key] = tr
	}
}
