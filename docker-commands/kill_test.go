package commands

import (
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/docker/docker/client"
)

func TestKill_DoesNotExitProcessOnContainerKillError(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestKill_DoesNotExitProcessOnContainerKillError_Helper")
	cmd.Env = append(os.Environ(), "KILL_HELPER_PROCESS=1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Kill should not terminate process on error. err=%v output=%s", err, string(output))
	}
}

func TestKill_DoesNotExitProcessOnContainerKillError_Helper(t *testing.T) {
	if os.Getenv("KILL_HELPER_PROCESS") != "1" {
		return
	}

	c, err := client.NewClientWithOpts(
		client.WithHost("tcp://127.0.0.1:1"),
		client.WithVersion("1.41"),
		client.WithHTTPClient(&http.Client{Timeout: 200 * time.Millisecond}),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	dockerClient := Client{c: c}
	dockerClient.Kill("invalid-container-id")
}
