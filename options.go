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

func AddTransformer(key string, fn Trasformer) Option {
	return func(f *Flatter) {
		f.trasformers[key] = fn
	}
}
