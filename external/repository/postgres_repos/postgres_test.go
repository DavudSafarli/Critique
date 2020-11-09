package postgres_repos

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net"
	"net/url"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/DavudSafarli/Critique/util"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type globalVars struct {
	storage *Storage
	storageErr error

	pool *dockertest.Pool
	resource *dockertest.Resource
	logWaiter docker.CloseWaiter
}
var vars = &globalVars{}
func TestMain(m *testing.M) {
	vars.initStorage()
	defer vars.cleanup()
	code := m.Run()
	os.Exit(code)
}
func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
func (v *globalVars) getPGScheme() *url.URL {
	return &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("myuser", "mypass"),
		Path:     "mydatabase",
		RawQuery: "sslmode=disable",
	}
}
func (v *globalVars) getContainerEnvVars() []string {
	scheme := v.getPGScheme()
	pw, _ := scheme.User.Password()
	return []string{
		"POSTGRES_USER=" + scheme.User.Username(),
		"POSTGRES_PASSWORD=" + pw,
		"POSTGRES_DB=" + scheme.Path,
	}
}
func (v *globalVars) getConnStr() string {
	scheme := v.getPGScheme()
	scheme.Host = fmt.Sprintf("localhost:%s", v.resource.GetPort("5432/tcp"))

	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		scheme.Host = net.JoinHostPort(v.resource.GetBoundIP("5432/tcp"), v.resource.GetPort("5432/tcp"))
	}
	return scheme.String()
}
// initStorage runs a new docker container, runs migrations and returns connection string of that postgres instance
func (v *globalVars) initStorage() {
	v.RunPostgresDocker()

	connStr := v.getConnStr()
	v.storage, v.storageErr = NewDbConnection(connStr)
	if v.storageErr != nil {
		panic(errors.Wrap(v.storageErr,"Failed to create a new Storage"))
	}
	// migrate
	migrationsFolder := "file://D:/GO/GOPATH/go/src/github.com/DavudSafarli/Critique/scripts/db/migrations"
	m, err := migrate.New(migrationsFolder, connStr)
	if err != nil {
		panic(err)
	}
	// Migrate all the way up ...
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	defer m.Close()
}
func (v *globalVars) cleanup() {
	err := v.pool.Purge(v.resource)
	panicOnErr(err)
	err = v.logWaiter.Close()
	panicOnErr(err)
	err = v.logWaiter.Wait()
	panicOnErr(err)
}


// RunPostgresDockerAndGetConnectionString does what its name says
func (v *globalVars) RunPostgresDocker() {
	var err error

	v.pool, err = dockertest.NewPool(util.GetDockerHost())
	panicOnErr(err)

	v.resource, err = v.pool.Run("postgres", "latest", v.getContainerEnvVars())
	panicOnErr(err)

	v.logWaiter, err = v.pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container: v.resource.Container.ID,
		//OutputStream: log.Writer(),
		ErrorStream: log.Writer(),
		Stderr:      true,
		//Stdout:       true,
		//Stream:       true,
	})
	panicOnErr(err)

	v.pool.MaxWait = 8 * time.Second
	err = v.pool.Retry(func() (err error) {
		a := v.getConnStr()
		db, e := sql.Open("pgx", a)
		if e != nil {
			return e
		}
		defer func() {
			cerr := db.Close()
			if err == nil {
				err = cerr
			}
		}()
		return db.Ping()
	})
	panicOnErr(err)
}

