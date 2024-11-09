package bgroups

import (
    "context"
    "errors"
    "github.com/Mides-Projects/Kyro/bgroups/model"
    "github.com/Mides-Projects/Operator/helper"
    "github.com/bytedance/sonic"
    "github.com/google/uuid"
    "github.com/nats-io/nats.go"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "strings"
    "sync"
)

type ServiceImpl struct {
    values map[string]*model.Group
    mu     sync.RWMutex

    ids   map[string]string
    idsMu sync.RWMutex

    col *mongo.Collection
    ctx context.Context
}

// cache caches the group information.
func (s *ServiceImpl) cache(g *model.Group) {
    s.mu.Lock()
    s.values[g.ID()] = g
    s.mu.Unlock()

    s.idsMu.Lock()
    s.ids[strings.ToLower(g.Name())] = g.ID()
    s.idsMu.Unlock()
}

// Values returns all the groups.
func (s *ServiceImpl) Values() []*model.Group {
    s.mu.RLock()
    defer s.mu.RUnlock()

    v := make([]*model.Group, 0, len(s.values))
    for _, g := range s.values {
        v = append(v, g)
    }

    return v
}

// LookupByID returns the group with the given ID.
func (s *ServiceImpl) LookupByID(id string) *model.Group {
    s.mu.RLock()
    defer s.mu.RUnlock()

    return s.values[id]
}

// LookupByName returns the group with the given name.
func (s *ServiceImpl) LookupByName(name string) *model.Group {
    s.idsMu.RLock()
    defer s.idsMu.RUnlock()

    if id, ok := s.ids[strings.ToLower(name)]; ok {
        return s.LookupByID(id)
    }

    return nil
}

// Insert inserts a new group with the given ID and name.
func (s *ServiceImpl) Insert(name string) (string, error) {
    if s.col == nil {
        return "", errors.New("no MongoDB collection")
    }

    g := model.NewGroup(uuid.New().String(), name)
    s.cache(g)

    go func() {
        // Insert the group into the MongoDB collection.
        if _, err := s.col.InsertOne(s.ctx, g); err != nil {
            panic(err)
        }

        helper.PublishNats(
            SubjectCreateGroup,
            map[string]interface{}{
                "id":   g.ID(),
                "name": g.Name(),
            },
        )
    }()

    helper.Log.Info("Successfully created group", "id", g.ID(), "name", name)

    return g.ID(), nil
}

// Hook initializes the group service.
func (s *ServiceImpl) Hook() error {
    if s.col != nil {
        return errors.New("collection already set")
    } else if helper.NatsClient == nil {
        return errors.New("nats client not set")
    }

    s.col = helper.MongoClient.Database("kyro").Collection("groups")
    // caching the context helps a lot with performance and memory usage
    s.ctx = context.Background()

    cur, err := s.col.Find(s.ctx, bson.M{})
    if err != nil {
        return err
    }

    for cur.Next(s.ctx) {
        var body map[string]interface{}
        g := &model.Group{}

        if err = cur.Decode(&body); err != nil {
            helper.Log.Error("failed to decode group", "error", err)
        } else if err = g.Unmarshal(body); err != nil {
            helper.Log.Error("failed to unmarshal group", "error", err, "body", body)
        } else {
            s.cache(g)
        }
    }

    helper.Log.Info("Successfully loaded " + string(len(s.values)) + " group(s) from the database!")

    if _, err := helper.NatsClient.Subscribe(SubjectCreateGroup, s.natsCreateGroup); err != nil {
        return errors.Join(errors.New("failed to subscribe to create group"), err)
    }

    return nil
}

// Service provides group management.
func (s *ServiceImpl) natsCreateGroup(msg *nats.Msg) {
    var body map[string]interface{}
    if err := sonic.Unmarshal(msg.Data, &body); err != nil {
        helper.Log.Error("failed to unmarshal create group message", "error", err)
    } else if id, ok := body["id"].(string); !ok {
        helper.Log.Error("create group message missing ID")
    } else if name, ok := body["name"].(string); !ok {
        helper.Log.Error("create group message missing name")
    } else {
        s.cache(model.NewGroup(id, name))

        helper.Log.Info("Successfully created group", "id", id, "name", name)
    }
}

func Service() *ServiceImpl {
    return service
}

var service = &ServiceImpl{
    values: make(map[string]*model.Group),
    ids:    make(map[string]string),
}

var SubjectCreateGroup = "kyro:create_group"
