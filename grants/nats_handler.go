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
    if service == nil {
        return errors.New("Kyro: no service")
    }

    service.ttlSet.Invalidate(id)

    return nil
}

// HandleQuit handles a quit request.
func (NatsHandler) HandleQuit(id string) error {
    if service == nil {
        return errors.New("Kyro: no service")
    } else if pi := service.Lookup(id); pi == nil {
        return errors.New("Kyro: player not found")
    } else {
        service.mu.Lock()
        delete(service.trackers, pi.ID())
        service.mu.Unlock()

        service.ttlSet.Invalidate(pi.ID())
    }

    return nil
}
