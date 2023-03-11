package flatr

type Option func(f *Flatter)

func WithPrefix(prefix string) Option {
	return func(f *Flatter) {
		f.prefix = prefix
	}
}

func WithSeparator(separator string) Option {
	return func(f *Flatter) {
		f.separator = separator
	}
}

func AddTransformer(tr Transformer) Option {
	return func(f *Flatter) {
		f.transformers = append(f.transformers, tr)
	}
}

func AddScopedTransformer(key string, tr Transformer) Option {
	return func(f *Flatter) {
		f.scopedTrasformers[key] = tr
	}
}
