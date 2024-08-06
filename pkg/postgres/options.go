package postgres

type Option func(postgres *postgres)

func MaxPoolSize(size int) Option {
	return func(postgres *postgres) {
		postgres.maxPoolSize = size
	}
}
