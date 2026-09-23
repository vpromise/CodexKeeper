package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/repository"
	repodto "cpa-usage-keeper/internal/repository/dto"
	"cpa-usage-keeper/internal/service"
)

func TestUsageIngestionKeepsOnlyNativeProviders(t *testing.T) {
	db := openUsageServiceTestDatabase(t)
	now := time.Now()
	payloads := []string{
		`{"provider":"codex","executor_type":"CodexExecutor","request_id":"native-codex","model":"pro01/gpt-test","tokens":{"input_tokens":20,"output_tokens":5,"total_tokens":25}}`,
		`{"provider":"claude","executor_type":"ClaudeExecutor","request_id":"native-claude","model":"claude-test","tokens":{"input_tokens":20,"output_tokens":5,"total_tokens":25}}`,
		`{"provider":"antigravity","executor_type":"ClaudeExecutor","request_id":"not-native-claude","model":"claude-test"}`,
		`{"provider":"openai","executor_type":"OpenAICompatExecutor","request_id":"not-native-codex","model":"gpt-test"}`,
		`{"provider":"gemini","request_id":"unsupported"}`,
		`{"request_id":"unknown-provider","model":"claude-test"}`,
	}
	rows := make([]repodto.RedisInboxInsert, 0, len(payloads))
	for _, payload := range payloads {
		rows = append(rows, repodto.RedisInboxInsert{Source: "usage", RawMessage: payload, PoppedAt: now})
	}
	if _, err := repository.InsertRedisUsageInboxMessages(db, rows); err != nil {
		t.Fatal(err)
	}
	s := service.NewSyncServiceWithOptions(db, service.SyncServiceOptions{BaseURL: "http://localhost"})
	result, err := s.ProcessRedisUsageInbox(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.InsertedEvents != 2 {
		t.Fatalf("expected two native events: %+v", result)
	}
	var events []entities.UsageEvent
	if err := db.Order("id").Find(&events).Error; err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 || events[0].Provider != "codex" || events[1].Provider != "claude" {
		t.Fatalf("provider scope leaked: %+v", events)
	}
	pending, err := repository.ListProcessableRedisUsageInbox(db, 100)
	if err != nil || len(pending) != 0 {
		t.Fatalf("unsupported messages remain retryable: %v %v", pending, err)
	}
}

func TestErrorEventsKeepOnlyNativeProviders(t *testing.T) {
	db := openUsageServiceTestDatabase(t)
	s := service.NewErrorEventService(db)
	for _, provider := range []string{" Codex ", "CLAUDE", "openai", "gemini", ""} {
		payload := fmt.Sprintf(`{"timestamp":"2026-09-23T14:00:00Z","provider":%q,"auth_index":"preview","status_code":429,"body":"rate limited"}`, provider)
		if err := s.StoreErrorEvent(context.Background(), payload, time.Now()); err != nil {
			t.Fatal(err)
		}
	}
	var rows []entities.ErrorEvent
	if err := db.Order("id").Find(&rows).Error; err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].Provider != "codex" || rows[1].Provider != "claude" {
		t.Fatalf("error event provider scope: %+v", rows)
	}
}
