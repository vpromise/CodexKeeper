package test

import (
	"context"

	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/service"
)

type analysisKeyStub struct {
	service.CPAAPIKeyProvider
	row       entities.CPAAPIKey
	listCalls int
}

func (s *analysisKeyStub) ListCPAAPIKeys(context.Context) ([]entities.CPAAPIKey, error) {
	s.listCalls++
	return []entities.CPAAPIKey{s.row, {ID: 99, APIKey: "sk-other654321", KeyAlias: "Other Key"}}, nil
}

func (s *analysisKeyStub) FindActiveCPAAPIKeyByID(_ context.Context, id int64) (entities.CPAAPIKey, error) {
	if id != s.row.ID {
		return entities.CPAAPIKey{}, service.ErrInvalidID
	}
	return s.row, nil
}
