package currency

import (
	"context"
	"time"

	"github.com/nburunova/currency-loader/src/services/log"
	"github.com/pkg/errors"
)

type ValutesLoader interface {
	Load(ctx context.Context) (Valutes, error)
}

type ValutesStorage interface {
	Save(ctx context.Context, valutes Valutes) error
}

type Service struct {
	loader  ValutesLoader
	storage ValutesStorage
	logger  log.Logger
}

func NewService(loader ValutesLoader, storage ValutesStorage, logger log.Logger) *Service {
	return &Service{
		loader:  loader,
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) update(ctx context.Context) error {
	s.logger.Debug("update currencies")
	newValutes, err := s.loader.Load(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot load valutes")
	}
	err = s.storage.Save(ctx, newValutes)
	if err != nil {
		return errors.Wrap(err, "cannot save valutes")
	}
	return nil
}

func (s *Service) Start(ctx context.Context, updatePeriod time.Duration) error {
	err := s.update(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot update currencies at first time")
	}
	ticker := time.NewTicker(updatePeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := s.update(ctx)
				if err != nil {
					s.logger.Error(errors.Wrap(err, "cannot update currencies"))
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}
