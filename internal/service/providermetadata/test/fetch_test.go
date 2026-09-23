package providermetadata_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestFetchFoldsReverseCompletionInRegistryOrder(t *testing.T) {
	fetcher := newGatedProviderFetcher(t)
	resultCh := startProviderFetch(context.Background(), fetcher)

	waitForSources(t, fetcher.entered, registrySourceOrder)
	reverseOrder := []string{"claude", "codex"}
	// 逐个释放并确认完成，确保真实完成顺序可控。
	for _, source := range reverseOrder {
		fetcher.release(source)
		waitForSources(t, fetcher.done, []string{source})
	}
	outcome := waitForFetchOutcome(t, resultCh)
	if outcome.err != nil {
		t.Fatalf("Fetch returned error: %v", outcome.err)
	}
	if !reflect.DeepEqual(outcome.snapshot.FetchedProviderTypes, registrySourceOrder) {
		t.Fatalf("FetchedProviderTypes = %#v", outcome.snapshot.FetchedProviderTypes)
	}
	gotAuthIndexes := make([]string, 0, len(outcome.snapshot.Credentials))
	for _, credential := range outcome.snapshot.Credentials {
		gotAuthIndexes = append(gotAuthIndexes, credential.AuthIndex)
	}
	wantAuthIndexes := []string{"codex-auth", "claude-auth"}
	if !reflect.DeepEqual(gotAuthIndexes, wantAuthIndexes) {
		t.Fatalf("auth indexes = %#v, want %#v", gotAuthIndexes, wantAuthIndexes)
	}
}

func TestFetchKeepsOtherSourcesWhenOneProviderFails(t *testing.T) {
	fetcher := newGatedProviderFetcher(t)
	fetcher.errors["claude"] = errors.New("claude unavailable")
	resultCh := startProviderFetch(context.Background(), fetcher)

	// 七个 endpoint 必须在任一结果返回前全部进入。
	waitForSources(t, fetcher.entered, registrySourceOrder)
	fetcher.releaseAll()
	waitForSources(t, fetcher.done, registrySourceOrder)
	outcome := waitForFetchOutcome(t, resultCh)
	if outcome.err == nil || outcome.err.Error() != "fetch claude api keys: claude unavailable" {
		t.Fatalf("error = %v", outcome.err)
	}
	wantTypes := []string{"codex"}
	if !reflect.DeepEqual(outcome.snapshot.FetchedProviderTypes, wantTypes) {
		t.Fatalf("FetchedProviderTypes = %#v, want %#v", outcome.snapshot.FetchedProviderTypes, wantTypes)
	}
	if len(outcome.snapshot.Credentials) != 1 {
		t.Fatalf("Credentials = %#v", outcome.snapshot.Credentials)
	}
}

func TestFetchPreservesCompletedSourcesAndWaitsForCancellation(t *testing.T) {
	fetcher := newGatedProviderFetcher(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resultCh := startProviderFetch(ctx, fetcher)

	waitForSources(t, fetcher.entered, registrySourceOrder)
	fetcher.release("codex")
	// 等待 Codex 报告完成，固定部分成功边界。
	waitForSources(t, fetcher.done, []string{"codex"})
	cancel()
	// 剩余六个 endpoint 必须全部退出，不允许 goroutine 泄漏。
	waitForSources(t, fetcher.done, []string{"claude"})
	outcome := waitForFetchOutcome(t, resultCh)
	if !reflect.DeepEqual(outcome.snapshot.FetchedProviderTypes, []string{"codex"}) || len(outcome.snapshot.Credentials) != 1 || outcome.snapshot.Credentials[0].AuthIndex != "codex-auth" {
		t.Fatalf("snapshot = %#v", outcome.snapshot)
	}
	wantError := "fetch claude api keys: context canceled"
	if outcome.err == nil || outcome.err.Error() != wantError {
		t.Fatalf("error = %v, want %q", outcome.err, wantError)
	}
}
