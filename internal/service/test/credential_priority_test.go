package test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/service"
	"gorm.io/gorm"
)

type priorityClientStub struct {
	patchName   string
	patchStatus int
	onPatch     func()
}

func (s *priorityClientStub) UpdateAuthFilePriority(_ context.Context, name string, priority int) (int, error) {
	s.patchName = name
	if s.patchStatus != 0 {
		return s.patchStatus, errors.New("upstream rejected")
	}
	if s.onPatch != nil {
		s.onPatch()
	}
	return http.StatusOK, nil
}

func loadPriority(t *testing.T, db *gorm.DB, authType entities.UsageIdentityAuthType, identity string) *int {
	t.Helper()
	var row entities.UsageIdentity
	if err := db.Where("auth_type = ? AND identity = ?", authType, identity).First(&row).Error; err != nil {
		t.Fatal(err)
	}
	return row.Priority
}

func requirePriority(t *testing.T, got *int, want int) {
	t.Helper()
	if got == nil || *got != want {
		t.Fatalf("priority=%v, want %d", got, want)
	}
}

func TestAuthFilePriorityPersistsExplicitZero(t *testing.T) {
	db := openMetadataTestDatabase(t, "priority-auth-file-zero.db")
	seedAuthFileCredential(t, db, "auth.json", "auth-index")
	client := &priorityClientStub{}
	refresher := &credentialStatusRefresherStub{}
	provider := service.NewCredentialPriorityService(db, client, refresher, &service.CredentialMutationLocks{})
	result, err := provider.SetAuthFilePriority(context.Background(), "auth-index", 0)
	if err != nil || result.Priority != 0 || client.patchName != "auth.json" {
		t.Fatalf("result=%+v err=%v patch=%q", result, err, client.patchName)
	}
	requirePriority(t, loadPriority(t, db, entities.UsageIdentityAuthTypeAuthFile, "auth-index"), 0)
	if refresher.count() != 1 {
		t.Fatalf("refresh count=%d", refresher.count())
	}
}

func TestAuthFilePriorityRejectsUpstreamConflict(t *testing.T) {
	db := openMetadataTestDatabase(t, "priority-auth-file-conflict.db")
	seedAuthFileCredential(t, db, "auth.json", "auth-index")
	client := &priorityClientStub{patchStatus: http.StatusConflict}
	provider := service.NewCredentialPriorityService(db, client, nil, &service.CredentialMutationLocks{})
	_, err := provider.SetAuthFilePriority(context.Background(), "auth-index", 3)
	if !errors.Is(err, service.ErrCredentialPriorityConflict) {
		t.Fatalf("err=%v, want conflict", err)
	}
	if got := loadPriority(t, db, entities.UsageIdentityAuthTypeAuthFile, "auth-index"); got != nil {
		t.Fatalf("unexpected local priority=%v", *got)
	}
}

func TestAuthFilePriorityPersistsSuccessfulPatch(t *testing.T) {
	db := openMetadataTestDatabase(t, "priority-auth-file-patch-success.db")
	seedAuthFileCredential(t, db, "auth.json", "auth-index")
	client := &priorityClientStub{}
	refresher := &credentialStatusRefresherStub{}
	result, err := service.NewCredentialPriorityService(db, client, refresher, &service.CredentialMutationLocks{}).SetAuthFilePriority(context.Background(), "auth-index", 3)
	if err != nil || result.Priority != 3 || refresher.count() != 1 {
		t.Fatalf("result=%+v err=%v refresh=%d", result, err, refresher.count())
	}
	requirePriority(t, loadPriority(t, db, entities.UsageIdentityAuthTypeAuthFile, "auth-index"), 3)
}

func TestPriorityRequestsRefreshIfLocalPersistenceFailsAfterCPASuccess(t *testing.T) {
	db := openMetadataTestDatabase(t, "priority-local-failure.db")
	seedAuthFileCredential(t, db, "auth.json", "target")
	refresher := &credentialStatusRefresherStub{}
	client := &priorityClientStub{}
	client.onPatch = func() {
		if err := db.Migrator().DropTable(&entities.UsageIdentity{}); err != nil {
			t.Fatal(err)
		}
	}
	_, err := service.NewCredentialPriorityService(db, client, refresher, &service.CredentialMutationLocks{}).SetAuthFilePriority(context.Background(), "target", 4)
	if err == nil || refresher.count() != 1 {
		t.Fatalf("err=%v refresh=%d", err, refresher.count())
	}
}

func TestPriorityUnsupportedTypeDoesNotPatch(t *testing.T) {
	db := openMetadataTestDatabase(t, "priority-unsupported.db")
	seedAuthFileCredential(t, db, "auth.json", "target")
	if err := db.Model(&entities.UsageIdentity{}).Where("identity = ?", "target").Update("type", "unknown").Error; err != nil {
		t.Fatal(err)
	}
	client := &priorityClientStub{}
	_, err := service.NewCredentialPriorityService(db, client, nil, &service.CredentialMutationLocks{}).SetAuthFilePriority(context.Background(), "target", 1)
	if !errors.Is(err, service.ErrCredentialPriorityUnsupported) {
		t.Fatalf("err=%v", err)
	}
	if client.patchName != "" {
		t.Fatalf("unexpected patch=%s", client.patchName)
	}
}

var _ service.CredentialPriorityClient = (*priorityClientStub)(nil)
