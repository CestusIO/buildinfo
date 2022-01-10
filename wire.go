//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package buildinfo

import "github.com/google/wire"

// ZapperLogProviderSet provides a zap logger
var BuildInfoProviderSet = wire.NewSet(
	ProvideBuildInfo,
)
