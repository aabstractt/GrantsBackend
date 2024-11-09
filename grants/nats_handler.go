package grants

import (
    "errors"
    "github.com/Mides-Projects/Zurita/handler"
)

type NatsHandler struct {
    handler.NatsHandler
}

// HandleHandshake handles a handshake request.
func (NatsHandler) HandleHandshake(id string, _ map[string]interface{}) error {
    if grantsService == nil {
        return errors.New("Kyro: no grantsService")
    }

    grantsService.ttlSet.Invalidate(id)

    return nil
}

// HandleQuit handles a quit request.
func (NatsHandler) HandleQuit(id string) error {
    if grantsService == nil {
        return errors.New("Kyro: no grantsService")
    } else if pi := grantsService.Lookup(id); pi == nil {
        return errors.New("Kyro: player not found")
    } else {
        grantsService.mu.Lock()
        delete(grantsService.trackers, pi.ID())
        grantsService.mu.Unlock()

        grantsService.ttlSet.Invalidate(pi.ID())
    }

    return nil
}
