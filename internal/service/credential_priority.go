package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/providers"
	"cpa-usage-keeper/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrCredentialPriorityValidation  = errors.New("credential priority request validation failed")
	ErrCredentialPriorityNotFound    = errors.New("credential priority target not found")
	ErrCredentialPriorityUnsupported = errors.New("credential priority is not supported")
	ErrCredentialPriorityConflict    = errors.New("credential priority target cannot be changed on its own")
)

type CredentialPriorityClient interface {
	UpdateAuthFilePriority(context.Context, string, int) (int, error)
}

type CredentialPriorityProvider interface {
	SetAuthFilePriority(context.Context, string, int) (CredentialPriorityResponse, error)
}

type CredentialPriorityResponse struct {
	AuthIndex string `json:"auth_index"`
	Priority  int    `json:"priority"`
}

type credentialPriorityService struct {
	db      *gorm.DB
	client  CredentialPriorityClient
	refresh MetadataRefresher
	locks   *CredentialMutationLocks
}

func NewCredentialPriorityService(db *gorm.DB, client CredentialPriorityClient, refresh MetadataRefresher, locks *CredentialMutationLocks) CredentialPriorityProvider {
	return &credentialPriorityService{db: db, client: client, refresh: refresh, locks: locks}
}

func (s *credentialPriorityService) SetAuthFilePriority(ctx context.Context, authIndex string, priority int) (CredentialPriorityResponse, error) {
	authIndex, err := s.validate(authIndex)
	if err != nil {
		return CredentialPriorityResponse{}, err
	}
	identity, err := s.findIdentity(ctx, entities.UsageIdentityAuthTypeAuthFile, authIndex)
	if err != nil {
		return CredentialPriorityResponse{}, err
	}
	if !providers.Supported(identity.Type) {
		return CredentialPriorityResponse{}, fmt.Errorf("%w: provider type %q", ErrCredentialPriorityUnsupported, identity.Type)
	}
	name := ""
	if identity.FileName != nil {
		name = strings.TrimSpace(*identity.FileName)
	}
	if name == "" {
		return CredentialPriorityResponse{}, fmt.Errorf("%w: auth file name is unavailable", ErrCredentialPriorityValidation)
	}

	defer s.locks.lockAuthFile(name)()

	statusCode, err := s.client.UpdateAuthFilePriority(ctx, name, priority)
	if err != nil {
		return CredentialPriorityResponse{}, priorityWriteError(statusCode, err)
	}
	defer s.requestRefresh()
	if err := repository.UpdateUsageIdentityPriority(ctx, s.db, entities.UsageIdentityAuthTypeAuthFile, authIndex, priority); err != nil {
		return CredentialPriorityResponse{}, fmt.Errorf("persist auth file priority: %w", err)
	}
	return CredentialPriorityResponse{AuthIndex: authIndex, Priority: priority}, nil
}

func (s *credentialPriorityService) validate(authIndex string) (string, error) {
	if s == nil || s.db == nil || s.client == nil {
		return "", fmt.Errorf("credential priority service is not configured")
	}
	authIndex = strings.TrimSpace(authIndex)
	if authIndex == "" {
		return "", fmt.Errorf("%w: auth_index is required", ErrCredentialPriorityValidation)
	}
	return authIndex, nil
}

func (s *credentialPriorityService) findIdentity(ctx context.Context, authType entities.UsageIdentityAuthType, authIndex string) (entities.UsageIdentity, error) {
	identity, err := repository.FindActiveUsageIdentityByAuthTypeAndIdentity(ctx, s.db, authType, authIndex)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.UsageIdentity{}, fmt.Errorf("%w: usage identity", ErrCredentialPriorityNotFound)
	}
	return identity, err
}

func (s *credentialPriorityService) requestRefresh() {
	if s.refresh != nil {
		s.refresh.RequestLocalMetadataRefresh()
	}
}

func priorityWriteError(statusCode int, err error) error {
	switch statusCode {
	case http.StatusNotFound:
		return fmt.Errorf("%w: upstream target", ErrCredentialPriorityNotFound)
	case http.StatusConflict:
		return fmt.Errorf("%w: upstream target", ErrCredentialPriorityConflict)
	default:
		return fmt.Errorf("update upstream priority: %w", err)
	}
}
