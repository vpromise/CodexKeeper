package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"cpa-usage-keeper/internal/auth"
	"cpa-usage-keeper/internal/timeutil"

	"github.com/gin-gonic/gin"
)

const (
	authSessionKindAdmin  = "admin"
	authSessionTimeLayout = "2006/01/02 15:04:05"
)

type authSessionListResponse struct {
	Items []authSessionItemResponse `json:"items"`
}

type authSessionItemResponse struct {
	ID          string    `json:"id"`
	Kind        string    `json:"kind"`
	Role        auth.Role `json:"role"`
	Source      string    `json:"source"`
	Current     bool      `json:"current,omitempty"`
	LoginAt     string    `json:"loginAt,omitempty"`
	LastSeenAt  string    `json:"lastSeenAt,omitempty"`
	ExpiresAt   string    `json:"expiresAt,omitempty"`
	LoginIP     string    `json:"loginIp,omitempty"`
	LastSeenIP  string    `json:"lastSeenIp,omitempty"`
	UserAgent   string    `json:"userAgent,omitempty"`
	Alias       *string   `json:"alias,omitempty"`
	sortSeenAt  time.Time `json:"-"`
	sortLoginAt time.Time `json:"-"`
}

func registerAuthSessionManagementRoutes(router gin.IRoutes, handler *authHandler) {
	router.GET("/auth/sessions", handler.listManagedSessions)
	router.PATCH("/auth/sessions/:id", handler.updateManagedSessionAlias)
	router.DELETE("/auth/sessions/:id", handler.revokeManagedSession)
}

func (h *authHandler) listManagedSessions(c *gin.Context) {
	if h == nil || !h.config.Enabled || h.sessions == nil {
		c.JSON(http.StatusOK, authSessionListResponse{Items: []authSessionItemResponse{}})
		return
	}
	records := h.sessions.List()
	c.JSON(http.StatusOK, authSessionListResponse{Items: buildAuthSessionItems(records, currentAuthSessionHash(c))})
}

func (h *authHandler) updateManagedSessionAlias(c *gin.Context) {
	if h == nil || !h.config.Enabled || h.sessions == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	sessionID := strings.TrimSpace(c.Param("id"))
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	alias, ok := parseUpdateAuthSessionAliasRequest(c)
	if !ok {
		return
	}
	if !h.sessions.UpdateAdminAliasByTokenHash(sessionID, alias) {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	items := buildAuthSessionItems(h.sessions.List(), currentAuthSessionHash(c))
	for _, item := range items {
		if item.ID == sessionID {
			c.JSON(http.StatusOK, item)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
}

func (h *authHandler) revokeManagedSession(c *gin.Context) {
	if h == nil || !h.config.Enabled || h.sessions == nil {
		c.Status(http.StatusNoContent)
		return
	}
	sessionID := strings.TrimSpace(c.Param("id"))
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	result := h.sessions.DeleteByTokenHash(sessionID)
	if result.Deleted == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	if sessionID == currentAuthSessionHash(c) {
		clearSessionCookie(c, h.config.BasePath, resolveSessionToken(c).CookieKind)
	}
	c.Status(http.StatusNoContent)
}

func buildAuthSessionItems(records []auth.SessionRecord, currentTokenHash string) []authSessionItemResponse {
	items := make([]authSessionItemResponse, 0, len(records))

	for _, record := range records {
		base := authSessionItemResponse{
			ID: record.TokenHash, Role: record.Role, Source: string(auth.NormalizeSessionSource(record.Source)),
			Current: record.TokenHash == currentTokenHash, LoginAt: formatAuthSessionTime(record.CreatedAt),
			LastSeenAt: formatAuthSessionTime(record.LastSeenAt), ExpiresAt: formatAuthSessionTime(record.ExpiresAt),
			LoginIP: record.LoginIP, LastSeenIP: record.LastSeenIP, UserAgent: record.UserAgent,
			sortSeenAt: record.LastSeenAt, sortLoginAt: record.CreatedAt,
		}
		if record.Role == auth.RoleAdmin {
			base.Kind = authSessionKindAdmin
			alias := record.Alias
			base.Alias = &alias
			items = append(items, base)
			continue
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Current != items[j].Current {
			return items[i].Current
		}
		if !items[i].sortSeenAt.Equal(items[j].sortSeenAt) {
			return items[i].sortSeenAt.After(items[j].sortSeenAt)
		}
		if !items[i].sortLoginAt.Equal(items[j].sortLoginAt) {
			return items[i].sortLoginAt.After(items[j].sortLoginAt)
		}
		if items[i].ID != items[j].ID {
			return items[i].ID < items[j].ID
		}
		return false
	})
	return items
}

func parseUpdateAuthSessionAliasRequest(c *gin.Context) (string, bool) {
	var payload map[string]json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return "", false
	}
	rawAlias, ok := payload["alias"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alias is required"})
		return "", false
	}
	if bytes.Equal(bytes.TrimSpace(rawAlias), []byte("null")) {
		return "", true
	}
	var alias string
	if err := json.Unmarshal(rawAlias, &alias); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alias must be a string or null"})
		return "", false
	}
	alias = strings.TrimSpace(alias)
	if err := validateUsageIdentityAlias(alias); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", false
	}
	return alias, true
}

func currentAuthSessionHash(c *gin.Context) string {
	if value, ok := c.Get(authTokenContextKey); ok {
		if token, ok := value.(string); ok && token != "" {
			return auth.SessionTokenHash(token)
		}
	}
	resolved := resolveSessionToken(c)
	if resolved.Token == "" {
		return ""
	}
	return auth.SessionTokenHash(resolved.Token)
}

func formatAuthSessionTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return timeutil.NormalizeStorageTime(value).Format(authSessionTimeLayout)
}
