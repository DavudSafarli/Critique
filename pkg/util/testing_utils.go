package util

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// GetDockerHost returns docker host uri based on OS
func GetDockerHost() string {
	if runtime.GOOS == "windows" {
		return "npipe:////./pipe/docker_engine"
	} else {
		return "unix:///var/run/docker.sock"
	}
}

// RunPostgresDockerAndGetConnectionString does what its name says
func RunPostgresDockerAndGetConnectionString(t *testing.T) string {
	t.Helper()

	pgURL := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("myuser", "mypass"),
		Path:     "mydatabase",
		RawQuery: "sslmode=disable",
	}

	pool, err := dockertest.NewPool(GetDockerHost())
	if err != nil {
		t.Fatalf("Could not connect to docker: %v", err)
	}

	pw, _ := pgURL.User.Password()
	env := []string{
		"POSTGRES_USER=" + pgURL.User.Username(),
		"POSTGRES_PASSWORD=" + pw,
		"POSTGRES_DB=" + pgURL.Path,
	}

	resource, err := pool.Run("postgres", "latest", env)
	if err != nil {
		t.Fatalf("Could not start postgres container: %v", err)
	}
	t.Cleanup(func() {
		err = pool.Purge(resource)
		if err != nil {
			t.Fatalf("Could not purge container: %v", err)
		}
	})

	pgURL.Host = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
	pgConnectionString := pgURL.String()

	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		pgURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}

	logWaiter, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container: resource.Container.ID,
		//OutputStream: log.Writer(),
		ErrorStream: log.Writer(),
		Stderr:      true,
		//Stdout:       true,
		//Stream:       true,
	})
	if err != nil {
		t.Fatalf("Could not connect to postgres container log output: %v", err)
	}

	t.Cleanup(func() {
		err = logWaiter.Close()
		if err != nil {
			t.Fatalf("Could not close container log: %v", err)
		}
		err = logWaiter.Wait()
		if err != nil {
			t.Fatalf("Could not wait for container log to close: %v", err)
		}
	})

	pool.MaxWait = 8 * time.Second
	err = pool.Retry(func() (err error) {
		db, err := sql.Open("pgx", pgConnectionString)
		if err != nil {
			return err
		}
		defer func() {
			cerr := db.Close()
			if err == nil {
				err = cerr
			}
		}()
		return db.Ping()
	})
	if err != nil {
		t.Fatalf("Could not connect to postgres container: %v", err)
	}

	return pgConnectionString
}

// MigrateDatabase migrates database all the way up
func MigrateDatabase(t *testing.T, pgURL string) {
	migrationsFolder := "file://D:/GO/GOPATH/go/src/github.com/DavudSafarli/Critique/scripts/db/migrations"
	m, err := migrate.New(migrationsFolder, pgURL)
	if err != nil {
		t.Fatalf("Cannot run migrations %v", err)
	}
	// Migrate all the way up ...
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}
	defer m.Close()
	t.Log("Successfully migrated the database")
}
