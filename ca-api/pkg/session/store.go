// pkg/session/store.go
package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dhruvpatel-10/signee/ca-api/db"
	"github.com/dhruvpatel-10/signee/ca-api/internal/domain/auth"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type SessionStore struct {
	cache        *cache.Cache
	userSessions map[uuid.UUID]map[string]*auth.Session
	mu           sync.RWMutex
	config       *Config
	metrics      *SessionMetrics
}

type Config struct {
	MaxSessions        int           // Maximum total sessions
	MaxUserSessions    int           // Maximum sessions per user
	SessionTTL         time.Duration // Default session TTL
	CleanupInterval    time.Duration // Cleanup interval
	InactivityTimeout  time.Duration // Inactivity timeout
	ExtendOnActivity   bool          // Auto-extend session on activity
	ConcurrentSessions bool          // Allow concurrent sessions per user
}

type SessionMetrics struct {
	TotalSessions   int64
	ActiveSessions  int64
	ExpiredSessions int64
	EvictedSessions int64
}

func NewSessionStore(config *Config) *SessionStore {
	if config == nil {
		config = &Config{
			MaxSessions:        10000,
			MaxUserSessions:    5,
			SessionTTL:         24 * time.Hour,
			CleanupInterval:    5 * time.Minute,
			InactivityTimeout:  30 * time.Minute,
			ExtendOnActivity:   true,
			ConcurrentSessions: true,
		}
	}

	store := &SessionStore{
		userSessions: make(map[uuid.UUID]map[string]*auth.Session),
		config:       config,
		metrics:      &SessionMetrics{},
	}

	// Setup cache with eviction hook
	store.cache = cache.New(config.SessionTTL, config.CleanupInterval)
	store.cache.OnEvicted(func(sessionID string, v interface{}) {
		if sess, ok := v.(*auth.Session); ok {
			store.mu.Lock()
			if userSessions, exists := store.userSessions[sess.UserID]; exists {
				delete(userSessions, sessionID)
				if len(userSessions) == 0 {
					delete(store.userSessions, sess.UserID)
				}
			}
			store.mu.Unlock()

			atomic.AddInt64(&store.metrics.ActiveSessions, -1)
			atomic.AddInt64(&store.metrics.ExpiredSessions, 1)
		}
	})

	return store
}

// Generate a cryptographically secure session ID
func generateSessionID() (string, error) {
	b := make([]byte, 32) // 256-bit random
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Create a new session
func (s *SessionStore) Create(ctx context.Context, user *db.User, metadata string) (*auth.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cache.ItemCount() >= s.config.MaxSessions {
		return nil, errors.New("session limit reached")
	}

	// Enforce per-user session policy
	userSessions := s.userSessions[user.ID]
	if userSessions == nil {
		userSessions = make(map[string]*auth.Session)
		s.userSessions[user.ID] = userSessions
	}

	if len(userSessions) >= s.config.MaxUserSessions {
		if !s.config.ConcurrentSessions {
			// Revoke all
			for id := range userSessions {
				s.cache.Delete(id)
			}
			userSessions = make(map[string]*auth.Session)
			s.userSessions[user.ID] = userSessions
		} else {
			// Remove one oldest session
			var oldestID string
			var oldestTime time.Time
			for id, sess := range userSessions {
				if oldestID == "" || sess.LastActivity.Before(oldestTime) {
					oldestID = id
					oldestTime = sess.LastActivity
				}
			}
			if oldestID != "" {
				s.cache.Delete(oldestID)
			}
		}
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}
	now := time.Now()

	session := &auth.Session{
		ID:           sessionID,
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Roles:        user.GetRoles(),
		Permissions:  user.GetPermissions(),
		IPAddress:    metadata.IPAddress,
		UserAgent:    metadata.UserAgent,
		CreatedAt:    now,
		LastActivity: now,
		MFAVerified:  metadata.MFAVerified,
	}

	s.cache.Set(sessionID, session, s.config.SessionTTL)
	userSessions[sessionID] = session

	atomic.AddInt64(&s.metrics.TotalSessions, 1)
	atomic.AddInt64(&s.metrics.ActiveSessions, 1)

	return session, nil
}

// Get session by ID
func (s *SessionStore) Get(ctx context.Context, sessionID string) (*auth.Session, error) {
	s.mu.RLock()
	cached, found := s.cache.Get(sessionID)
	s.mu.RUnlock()

	if !found {
		return nil, ErrSessionNotFound
	}
	session := cached.(*auth.Session)

	// Check inactivity timeout
	if time.Since(session.LastActivity) > s.config.InactivityTimeout {
		s.Remove(ctx, sessionID)
		return nil, ErrSessionInactive
	}

	// Update activity inline
	if s.config.ExtendOnActivity {
		s.mu.Lock()
		session.LastActivity = time.Now()
		s.cache.Set(sessionID, session, s.config.SessionTTL) // extend TTL
		s.mu.Unlock()
	}

	return session, nil
}

// Remove session
func (s *SessionStore) Remove(ctx context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.cache.Get(sessionID); !found {
		return ErrSessionNotFound
	}
	s.cache.Delete(sessionID)
	atomic.AddInt64(&s.metrics.EvictedSessions, 1)
	return nil
}
