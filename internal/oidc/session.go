// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oidc

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/gomodule/redigo/redis"
)

// ErrSessionExpired is returned when a session has exceeded idle timeout.
var ErrSessionExpired = errors.New("session idle timeout exceeded")

// Flow represents an OAuth2 flow state stored in Redis.
type Flow struct {
	State     string    `json:"state"`
	Nonce     string    `json:"nonce"`
	Verifier  string    `json:"verifier"`
	CreatedAt time.Time `json:"created_at"`
}

// Session represents a user session stored in Redis.
type Session struct {
	Sub       string    `json:"sub"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Groups    []string  `json:"groups"`
	IDToken   string    `json:"id_token"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}

// GenerateID generates a cryptographically secure random ID.
func GenerateID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashID returns SHA-256 hash of ID for storage (so raw IDs aren't exposed in Redis dumps).
func HashID(id string) string {
	hash := sha256.Sum256([]byte(id))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// StoreFlow stores an OAuth2 flow in Redis with 10 min TTL.
func StoreFlow(conn redis.Conn, flowID string, flow *Flow) error {
	data, err := json.Marshal(flow)
	if err != nil {
		return fmt.Errorf("failed to marshal flow: %w", err)
	}

	key := fmt.Sprintf("fsui:oidc:flow:%s", HashID(flowID))
	_, err = conn.Do("SETEX", key, 600, string(data))
	if err != nil {
		return fmt.Errorf("failed to store flow: %w", err)
	}
	return nil
}

// GetFlow retrieves and deletes a flow from Redis (one-time use).
func GetFlow(conn redis.Conn, flowID string) (*Flow, error) {
	key := fmt.Sprintf("fsui:oidc:flow:%s", HashID(flowID))
	data, err := redis.Bytes(conn.Do("GETDEL", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil, fmt.Errorf("flow not found")
		}
		return nil, fmt.Errorf("failed to get flow: %w", err)
	}

	var flow Flow
	err = json.Unmarshal(data, &flow)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal flow: %w", err)
	}
	return &flow, nil
}

// StoreSession stores a user session in Redis.
func StoreSession(conn redis.Conn, sessionID string, session *Session) error {
	config := configuration.GetConfiguration()

	// Update LastSeen timestamp
	if session.LastSeen.IsZero() {
		session.LastSeen = time.Now()
	}

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	key := fmt.Sprintf("fsui:session:%s", HashID(sessionID))

	// Calculate TTL as minimum of idle timeout and absolute session TTL
	absoluteTTL := config.SessionTTL
	idleTTL := config.SessionIdleTimeout

	ttl := absoluteTTL
	if idleTTL > 0 && idleTTL < absoluteTTL {
		ttl = idleTTL
	}

	_, err = conn.Do("SETEX", key, ttl, string(data))
	if err != nil {
		return fmt.Errorf("failed to store session: %w", err)
	}
	return nil
}

// GetSession retrieves a session from Redis.
func GetSession(conn redis.Conn, sessionID string) (*Session, error) {
	key := fmt.Sprintf("fsui:session:%s", HashID(sessionID))
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session Session
	err = json.Unmarshal(data, &session)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}
	return &session, nil
}

// RefreshSessionTTL updates the idle timeout for a session.
// Refreshes when time.Since(session.LastSeen) >= min(60s, idle/4).
// Returns ErrSessionExpired if the session has exceeded idle timeout.
// Returns other errors for transient failures (e.g., Redis errors).
func RefreshSessionTTL(conn redis.Conn, sessionID string) error {
	config := configuration.GetConfiguration()
	key := fmt.Sprintf("fsui:session:%s", HashID(sessionID))

	// Get current session data
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			// Session not found (may have expired in Redis)
			return ErrSessionExpired
		}
		// Transient Redis error: return without deleting
		return fmt.Errorf("failed to get session: %w", err)
	}

	var session Session
	err = json.Unmarshal(data, &session)
	if err != nil {
		return fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Check idle timeout
	idleTTL := config.SessionIdleTimeout
	absoluteTTL := config.SessionTTL

	if idleTTL > 0 && time.Since(session.LastSeen) > time.Duration(idleTTL)*time.Second {
		return ErrSessionExpired
	}

	// Determine when to refresh: min(60s, idle/4)
	refreshThreshold := 60 * time.Second
	if idleTTL > 0 {
		refreshThreshold = time.Duration(idleTTL/4) * time.Second
		if refreshThreshold > 60*time.Second {
			refreshThreshold = 60 * time.Second
		}
	}

	// Only refresh if enough time has passed since LastSeen
	if time.Since(session.LastSeen) >= refreshThreshold {
		// Update LastSeen
		session.LastSeen = time.Now()
		updatedData, err := json.Marshal(session)
		if err != nil {
			return fmt.Errorf("failed to marshal session: %w", err)
		}

		// Calculate new TTL: min(idle, remaining absolute = absoluteTTL - time.Since(CreatedAt)), at least 1s
		newTTL := absoluteTTL
		if idleTTL > 0 && idleTTL < newTTL {
			newTTL = idleTTL
		}

		// Ensure remaining absolute TTL doesn't exceed absolute session TTL
		elapsedFromCreation := time.Since(session.CreatedAt)
		remainingAbsolute := time.Duration(absoluteTTL)*time.Second - elapsedFromCreation
		if remainingAbsolute < time.Duration(newTTL)*time.Second {
			newTTL = int(remainingAbsolute.Seconds())
		}

		// Ensure at least 1s
		if newTTL < 1 {
			newTTL = 1
		}

		// Use SET key val XX EX to atomically update only if key exists
		// This prevents resurrection of a concurrently deleted session
		_, err = conn.Do("SET", key, string(updatedData), "XX", "EX", newTTL)
		if err != nil {
			return fmt.Errorf("failed to refresh session: %w", err)
		}
	}

	return nil
}

// DeleteSession deletes a session from Redis.
func DeleteSession(conn redis.Conn, sessionID string) error {
	key := fmt.Sprintf("fsui:session:%s", HashID(sessionID))
	_, err := conn.Do("DEL", key)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
