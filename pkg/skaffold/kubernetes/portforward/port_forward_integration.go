package portforward

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
)

//This is testing a port forward + stop + restart in a simulated dev cycle
func WhiteBox_PortForwardCycle(namespace string, t *testing.T) {
	em := NewEntryManager(os.Stdout)
	portForwardEvent = func(entry *portForwardEntry) {}
	ctx := context.Background()
	localPort := retrieveAvailablePort(9000, em.forwardedPorts)
	pfe := &portForwardEntry{
		resource: latest.PortForwardResource{
			Type:      "deployment",
			Name:      "leeroy-web",
			Namespace: namespace,
			Port:      8080,
		},
		containerName: "dummy container",
		localPort:     localPort,
	}

	defer em.Stop()
	if err := em.forwardPortForwardEntry(ctx, pfe); err != nil {
		t.Fatalf("failed to forward port: %s", err)
	}
	em.Stop()

	time.Sleep(2 * time.Second)

	logrus.Info("getting next port...")
	nextPort := retrieveAvailablePort(localPort, em.forwardedPorts)

	// theoretically we should be able to bind to the very same port
	// this might get flaky when multiple tests are ran. However
	// we shouldn't collide with our own process because of poor cleanup
	if nextPort != localPort {
		t.Fatalf("the same port should be still open, instead first port: %d != second port: %d", localPort, nextPort)
	}

	defer em.Stop()
	if err := em.forwardPortForwardEntry(ctx, pfe); err != nil {
		t.Fatalf("failed to forward port: %s", err)
	}

}
