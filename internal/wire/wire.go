//go:build wireinject
// +build wireinject

package wire

import (
	"dove/internal/app"

	"github.com/google/wire"
)

// InitializeApp 初始化应用
func InitializeApp() (*app.App, error) {
	wire.Build(ProviderSet)
	return &app.App{}, nil
}
