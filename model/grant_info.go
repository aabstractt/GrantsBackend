package model

import (
	"errors"
	"time"
)

type GrantInfo struct {
	id string

	grant Grant

	addedBy string
	addedAt time.Time

	expiresAt time.Time

	revokedBy string
	revokedAt time.Time

	scopes []string
}

// ID returns the ID of the grant.
func (gi *GrantInfo) ID() string {
	return gi.id
}

// Grant returns the grant of the grant.
func (gi *GrantInfo) Grant() Grant {
	return gi.grant
}

// AddedBy returns the added by of the grant.
func (gi *GrantInfo) AddedBy() string {
	return gi.addedBy
}

// AddedAt returns the added at of the grant.
func (gi *GrantInfo) AddedAt() time.Time {
	return gi.addedAt
}

// ExpiresAt returns the expires at of the grant.
func (gi *GrantInfo) ExpiresAt() time.Time {
	return gi.expiresAt
}

// RevokedBy returns the revoked by of the grant.
func (gi *GrantInfo) RevokedBy() string {
	return gi.revokedBy
}

// SetRevokedBy sets the revoked by of the grant.
func (gi *GrantInfo) SetRevokedBy(revokedBy string) {
	gi.revokedBy = revokedBy
}

// RevokedAt returns the revoked at of the grant.
func (gi *GrantInfo) RevokedAt() time.Time {
	return gi.revokedAt
}

// SetRevokedAt sets the revoked at of the grant.
func (gi *GrantInfo) SetRevokedAt(revokedAt time.Time) {
	gi.revokedAt = revokedAt
}

// Expired returns if the grant is expired.
func (gi *GrantInfo) Expired() bool {
	if gi.revokedAt.Unix() != 0 {
		return true
	}

	return gi.expiresAt.Unix() > 0 && time.Now().After(gi.expiresAt)
}

// Scopes returns the scopes of the grant.
func (gi *GrantInfo) Scopes() []string {
	return gi.scopes
}

// SetScopes sets the scopes of the grant.
func (gi *GrantInfo) SetScopes(scopes []string) {
	gi.scopes = scopes
}

// Marshal returns the grant info as a map.
func (gi *GrantInfo) Marshal() map[string]interface{} {
	body := map[string]interface{}{
		"_id":   gi.id,
		"grant": gi.grant.Marshal(),

		"added_by": gi.addedBy,
		"added_at": gi.addedAt.Unix(),

		"expires_at": gi.expiresAt.Unix(),
		"scopes":     gi.scopes,
	}

	if gi.revokedBy != "" && gi.revokedAt.Unix() != 0 {
		body["revoked_by"] = gi.revokedBy
		body["revoked_at"] = gi.revokedAt.Unix()
	}

	return body
}

// Unmarshal unmarshals the grant info from the given map.
func (gi *GrantInfo) Unmarshal(body map[string]interface{}) error {
	id, ok := body["_id"].(string)
	if !ok {
		return errors.New("_id is not a string")
	}
	gi.id = id

	grant := &Grant{}
	if err := grant.Unmarshal(body["grant"].(map[string]interface{})); err != nil {
		return err
	}
	gi.grant = *grant // reassign the pointer to the value

	addedBy, ok := body["added_by"].(string)
	if !ok {
		return errors.New("added_by is not a string")
	}
	gi.addedBy = addedBy

	addedAt, ok := body["added_at"].(int64)
	if !ok {
		return errors.New("added_at is not an integer")
	}
	gi.addedAt = time.Unix(addedAt, 0)

	expiresAt, ok := body["expires_at"].(int64)
	if !ok {
		return errors.New("expires_at is not an integer")
	}
	gi.expiresAt = time.Unix(expiresAt, 0)

	if revokedBy, ok := body["revoked_by"].(string); ok {
		gi.revokedBy = revokedBy
	}

	if revokedAt, ok := body["revoked_at"].(int64); ok {
		gi.revokedAt = time.Unix(revokedAt, 0)
	}

	if scopes, ok := body["scopes"].([]interface{}); !ok {
		return errors.New("scopes is not an array")
	} else {
		for _, scope := range scopes {
			if s, ok := scope.(string); !ok {
				return errors.New("scope is not a string")
			} else {
				gi.scopes = append(gi.scopes, s)
			}
		}
	}

	return nil
}
