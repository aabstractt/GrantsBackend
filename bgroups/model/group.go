package model

import (
    "errors"
    "slices"
    "sync"
)

type Group struct {
    id string

    name        string // Name is the name of the group.
    displayName string // DisplayName is the display name of the group.

    charColor string // CharColor is the color of the display name.

    prefix string // Prefix is the first part of the display name.
    suffix string // Suffix is the last part of the display name.

    chatPrefix string // ChatPrefix is the prefix for chat messages.
    chatSuffix string // ChatSuffix is the suffix for chat messages.

    permissionsMu sync.RWMutex // PermissionsMu is the mutex for the permissions.
    permissions   []string     // Permissions is the list of permissions the group has.
}

func NewGroup(id, name string) *Group {
    return &Group{
        id:   id,
        name: name,
    }
}

// ID returns the ID of the group.
func (g *Group) ID() string {
    return g.id
}

// Name returns the name of the group.
func (g *Group) Name() string {
    return g.name
}

// DisplayName returns the display name of the group.
func (g *Group) DisplayName() string {
    return g.displayName
}

// SetDisplayName sets the display name of the group.
func (g *Group) SetDisplayName(displayName string) {
    g.displayName = displayName
}

// CharColor returns the color of the display name.
func (g *Group) CharColor() string {
    return g.charColor
}

// SetCharColor sets the color of the display name.
func (g *Group) SetCharColor(charColor string) {
    g.charColor = charColor
}

// Prefix returns the prefix of the display name.
func (g *Group) Prefix() string {
    return g.prefix
}

// SetPrefix sets the prefix of the display name.
func (g *Group) SetPrefix(prefix string) {
    g.prefix = prefix
}

// Suffix returns the suffix of the display name.
func (g *Group) Suffix() string {
    return g.suffix
}

// SetSuffix sets the suffix of the display name.
func (g *Group) SetSuffix(suffix string) {
    g.suffix = suffix
}

// ChatPrefix returns the chat prefix of the group.
func (g *Group) ChatPrefix() string {
    return g.chatPrefix
}

// SetChatPrefix sets the chat prefix of the group.
func (g *Group) SetChatPrefix(chatPrefix string) {
    g.chatPrefix = chatPrefix
}

// ChatSuffix returns the chat suffix of the group.
func (g *Group) ChatSuffix() string {
    return g.chatSuffix
}

// SetChatSuffix sets the chat suffix of the group.
func (g *Group) SetChatSuffix(chatSuffix string) {
    g.chatSuffix = chatSuffix
}

// Permissions returns the permissions of the group.
func (g *Group) Permissions() []string {
    g.permissionsMu.RLock()
    defer g.permissionsMu.RUnlock()

    return g.permissions
}

// AddPermission adds the permission to the group.
func (g *Group) AddPermission(permission string) {
    g.permissionsMu.Lock()
    g.permissions = append(g.permissions, permission)
    g.permissionsMu.Unlock()
}

// RemovePermission removes the permission from the group.
func (g *Group) RemovePermission(permission string) {
    g.permissionsMu.Lock()

    if indx := slices.Index(g.permissions, permission); indx != -1 {
        g.permissions = append(g.permissions[:indx], g.permissions[indx+1:]...)
    }

    g.permissionsMu.Unlock()
}

// Marshal marshals the group into a map.
func (g *Group) Marshal() map[string]interface{} {
    body := map[string]interface{}{
        "_id":  g.id,
        "name": g.name,
    }
    if g.displayName != "" {
        body["display_name"] = g.displayName
    }

    if g.charColor != "" {
        body["char_color"] = g.charColor
    }

    if g.prefix != "" {
        body["prefix"] = g.prefix
    }

    if g.suffix != "" {
        body["suffix"] = g.suffix
    }

    if g.chatPrefix != "" {
        body["chat_prefix"] = g.chatPrefix
    }

    if g.chatSuffix != "" {
        body["chat_suffix"] = g.chatSuffix
    }

    if len(g.permissions) > 0 {
        body["permissions"] = g.Permissions()
    }

    return body
}

// Unmarshal unmarshals the body into the group.
func (g *Group) Unmarshal(body map[string]interface{}) error {
    id, ok := body["_id"].(string)
    if !ok {
        return errors.New("_id is not a string")
    }
    g.id = id

    name, ok := body["name"].(string)
    if !ok {
        return errors.New("name is not a string")
    }
    g.name = name

    if displayName, ok := body["display_name"].(string); ok {
        g.displayName = displayName
    }

    if charColor, ok := body["char_color"].(string); ok {
        g.charColor = charColor
    }

    if prefix, ok := body["prefix"].(string); ok {
        g.prefix = prefix
    }

    if suffix, ok := body["suffix"].(string); ok {
        g.suffix = suffix
    }

    if chatPrefix, ok := body["chat_prefix"].(string); ok {
        g.chatPrefix = chatPrefix
    }

    if chatSuffix, ok := body["chat_suffix"].(string); ok {
        g.chatSuffix = chatSuffix
    }

    if permissions, ok := body["permissions"].([]string); ok {
        g.permissionsMu.Lock()
        g.permissions = permissions
        g.permissionsMu.Unlock()
    }

    return nil
}
