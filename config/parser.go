package config

import "github.com/caarlos0/env/v11"

func parseEnv[T any](opts env.Options) (T, error) {
	var t T

	if err := env.Parse(&t); err != nil {
		return t, err
	}

	// override with PREFIX_XXX_XXX if when it has valu
	// this is optional no need handle error because if it not found it will use default value

	//nolint:all
	env.ParseWithOptions(&t, opts)

	return t, nil
}
