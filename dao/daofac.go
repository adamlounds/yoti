package dao

import (
	"github.com/adamlounds/yoti/store"
	"github.com/rs/zerolog"
)

type Daofac struct {
	logger zerolog.Logger
	Document *DocumentDao
}

func NewFactory(logger zerolog.Logger, storeFac *store.StoreFac) *Daofac {
	return &Daofac{logger, &DocumentDao{storeFac.DocumentStore}}
}

