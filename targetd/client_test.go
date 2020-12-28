package targetd

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
)

const (
	testFilePath = "/datafile-go-targetd-test"
)

var (
	testHost = "http://127.0.0.1:18700"
)

func setup() (*Client, error) {
	client, err := New(testHost, "testing", "secret_password", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to targetd.New: %w", err)
	}

	return client, nil
}

func TestMain(m *testing.M) {
	os.Exit(IntegrationTestRunner(m))
}

// IntegrationTestRunner is all integration test
func IntegrationTestRunner(m *testing.M) int {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	runOption := &dockertest.RunOptions{
		Name: "go-targetd",
		Mounts: []string{
			"/sys/kernel/config:/sys/kernel/config",
			"/lib/modules:/lib/modules",
			"/dev:/dev",
		},
		ExposedPorts: []string{
			"18700/tcp",
		},
		Privileged: true,
	}

	// Build and run the given Dockerfile
	resource, err := pool.BuildAndRunWithOptions("../testing/docker/Dockerfile", runOption)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		testHost = fmt.Sprintf("http://localhost:%s", resource.GetPort("18700/tcp"))

		if _, err := net.DialTimeout("tcp", testHost, 10*time.Second); err != nil {
			return fmt.Errorf("cloud not dial targetd host: %w", err)
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	defer func() {
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	//if out, err := exec.CommandContext(context.Background(), "../test/scripts/init.sh").CombinedOutput(); err != nil {
	//	log.Printf("init.sh return err: %+v (out: %+v)", err, out)
	//	return 1
	//}

	// setup filesystem
	if err := initializeFile(); err != nil {
		log.Fatalf("Cloud not initialize device: %+v", err)
	}

	code := m.Run()

	if err := initializeFile(); err != nil {
		log.Fatalf("Cloud not initalize device: %+v", err)
	}

	return code
}

func initializeFile() error {
	// TODO: execute in targetd container
	out, err := exec.Command("rm", "-f", testFilePath).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete %s (out: %s): %w", testFilePath, string(out), err)
	}

	return nil
}
