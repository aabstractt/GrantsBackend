package grants

import (
    "context"
    "errors"
    model2 "github.com/Mides-Projects/Kyro/grants/model"
    "github.com/Mides-Projects/Operator/helper"
    "github.com/Mides-Projects/Quark"
    "github.com/Mides-Projects/Zurita"
    pimodel "github.com/Mides-Projects/Zurita/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "sync"
    "time"
)

type ServiceImpl struct {
    trackers map[string]*model2.Tracker
    mu       sync.RWMutex

    ttlSet *Quark.Set
    // Player collection from MongoDB.
    col *mongo.Collection
    ctx context.Context
}

// cache caches the tracker information.
func (s *ServiceImpl) cache(t *model2.Tracker, keep bool) {
    s.mu.Lock()
    s.trackers[t.ID()] = t
    s.mu.Unlock()

    if keep {
        return
    }

    s.ttlSet.Set(t.ID())
}

// Lookup returns the tracker with the given ID.
// This method is thread-safe because it only reads the cache.
func (s *ServiceImpl) Lookup(id string) *model2.Tracker {
    s.mu.RLock()
    defer s.mu.RUnlock()

    return s.trackers[id]
}

// UnsafeLookup returns the tracker with the given ID
// first by checking the cache and then the MongoDB collection.
// This method is not thread-safe.
func (s *ServiceImpl) UnsafeLookup(id string) (*model2.Tracker, error) {
    if t := s.Lookup(id); t != nil {
        return t, nil
    } else if s.col == nil {
        return nil, errors.New("no MongoDB collection")
    } else if s.ctx == nil {
        return nil, errors.New("no context")
    }

    // Fetch the grants from the MongoDB collection.
    cur, err := s.col.Find(s.ctx, bson.M{"source_id": id})
    if err != nil {
        return nil, err
    }

    t := model2.NewTracker(id)
    for cur.Next(s.ctx) {
        var body map[string]interface{}
        if err = cur.Decode(&body); err != nil {
            return nil, err
        }

        gi := &model2.GrantInfo{}
        if err = gi.Unmarshal(body); err != nil {
            return nil, err
        }

        if gi.Expired() {
            t.AddExpired(*gi)
        } else {
            t.AddActive(gi)
        }
    }

    return t, nil
}

// HandleLookup handles the lookup of a player.
func (s *ServiceImpl) HandleLookup(id string, idSrc bool) (*model2.Tracker, error) {
    var (
        pi  *pimodel.PlayerInfo
        err error
    )
    if idSrc {
        pi, err = Zurita.Service().UnsafeLookupByID(id)
    } else {
        pi, err = Zurita.Service().UnsafeLookupByName(id)
    }

    if err != nil {
        return nil, err
    } else if pi == nil {
        return nil, nil
    } else if t, err := s.UnsafeLookup(pi.ID()); err != nil {
        return nil, err
    } else {
        if t == nil {
            t = model2.NewTracker(pi.ID())
        }

        // Cache the tracker if it does not exist.
        if s.Lookup(pi.ID()) == nil {
            s.cache(t, pi.Online())
        }

        return t, nil
    }
}

// Hook initializes the service.
func (s *ServiceImpl) Hook() error {
    if s.ttlSet != nil {
        return errors.New("GrantsX: TTL set already set")
    } else if s.col != nil {
        return errors.New("GrantsX: mongo collection already set")
    }

    // caching the context helps a lot with performance and memory usage
    s.ctx = context.Background()

    s.ttlSet = Quark.NewSet(
        1*time.Hour,
        1*time.Hour,
    )
    s.ttlSet.SetListener(func(id string, r Quark.Reason) {
        if r == Quark.ManualReason {
            return
        }

        pi := Zurita.Service().LookupByID(id)
        if pi != nil && pi.Online() {
            return // No clear the tracker if the player is online.
        }

        s.mu.Lock()
        delete(s.trackers, id)
        s.mu.Unlock()
    })

    s.col = helper.MongoClient.Database(helper.MongoDBName).Collection("grants")

    Zurita.Service().SetNatsHandler(NatsHandler{})

    return nil
}

// Service returns the service.
func Service() *ServiceImpl {
    return service
}

var service = &ServiceImpl{
    trackers: make(map[string]*model2.Tracker),
}
