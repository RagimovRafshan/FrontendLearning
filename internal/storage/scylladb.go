package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"messenger/config"
)

// ScyllaDBClient обертка над gocql.Session для работы со ScyllaDB
type ScyllaDBClient struct {
	session  *gocql.Session
	keyspace string
}

// NewScyllaDBClient создает новый клиент ScyllaDB
func NewScyllaDBClient(cfg config.ScyllaDBConfig) (*ScyllaDBClient, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)
	cluster.Keyspace = cfg.Keyspace
	cluster.Timeout = cfg.Timeout
	cluster.ConnectTimeout = cfg.ConnectTimeout
	cluster.Consistency = parseConsistency(cfg.Consistency)

	if cfg.Username != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

	// Настройка пула соединений
	cluster.NumConns = 4
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(
		gocql.RoundRobinHostPolicy(),
	)

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create ScyllaDB session: %w", err)
	}

	client := &ScyllaDBClient{
		session:  session,
		keyspace: cfg.Keyspace,
	}

	return client, nil
}

// parseConsistency преобразует строку в константу Consistency
func parseConsistency(consistency string) gocql.Consistency {
	switch consistency {
	case "ALL":
		return gocql.All
	case "ANY":
		return gocql.Any
	case "ONE":
		return gocql.One
	case "TWO":
		return gocql.Two
	case "THREE":
		return gocql.Three
	case "QUORUM":
		return gocql.Quorum
	case "LOCAL_ONE":
		return gocql.LocalOne
	case "LOCAL_QUORUM":
		return gocql.LocalQuorum
	case "EACH_QUORUM":
		return gocql.EachQuorum
	default:
		return gocql.LocalQuorum
	}
}

// Session возвращает базовую сессию gocql
func (c *ScyllaDBClient) Session() *gocql.Session {
	return c.session
}

// Keyspace возвращает имя keyspace
func (c *ScyllaDBClient) Keyspace() string {
	return c.keyspace
}

// Close закрывает соединение со ScyllaDB
func (c *ScyllaDBClient) Close() {
	if c.session != nil {
		c.session.Close()
	}
}

// HealthCheck проверяет подключение к ScyllaDB
func (c *ScyllaDBClient) HealthCheck(ctx context.Context) error {
	query := c.session.Query("SELECT now() FROM system.local LIMIT 1").WithContext(ctx)
	var result time.Time
	err := query.Scan(&result)
	if err != nil {
		return fmt.Errorf("ScyllaDB health check failed: %w", err)
	}
	return nil
}

// CreateKeyspace создает keyspace если он не существует
func (c *ScyllaDBClient) CreateKeyspace(ctx context.Context, name string, replicationFactor int) error {
	queryStr := fmt.Sprintf(
		`CREATE KEYSPACE IF NOT EXISTS %s 
		WITH REPLICATION = { 
			'class' : 'NetworkTopologyStrategy', 
			'datacenter1' : %d 
		}`,
		name,
		replicationFactor,
	)

	query := c.session.Query(queryStr).WithContext(ctx)
	if err := query.Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace: %w", err)
	}

	log.Printf("Keyspace '%s' created or already exists", name)
	return nil
}

// Execute выполняет CQL запрос
func (c *ScyllaDBClient) Execute(ctx context.Context, stmt string, values ...interface{}) error {
	query := c.session.Query(stmt, values...).WithContext(ctx)
	return query.Exec()
}

// Query выполняет CQL запрос и возвращает iterator
func (c *ScyllaDBClient) Query(ctx context.Context, stmt string, values ...interface{}) *gocql.Iter {
	query := c.session.Query(stmt, values...).WithContext(ctx)
	return query.Iter()
}
