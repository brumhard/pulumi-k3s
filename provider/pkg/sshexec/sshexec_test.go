package sshexec_test

import (
	"bytes"
	"log"
	"net"
	"testing"
	"time"

	_ "embed"

	"github.com/brumhard/pulumi-k3s/provider/pkg/sshexec"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/id_rsa
var sshkey []byte

//go:generate docker build -t ssh-test .
func Test_Client_Run(t *testing.T) {
	pool, err := dockertest.NewPool("")
	assert.NoError(t, err)

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "ssh-test",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	assert.NoError(t, err)

	if err := pool.Retry(func() error {
		_, err = net.DialTimeout("tcp", resource.GetHostPort("22/tcp"), 5*time.Second)
		return err
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	defer pool.Purge(resource)

	t.Run("client can execute commands", func(t *testing.T) {
		client, err := sshexec.NewClient(resource.GetHostPort("22/tcp"), "sshuser", sshkey)
		assert.NoError(t, err)

		stdout := &bytes.Buffer{}
		err = client.Run(&sshexec.Cmd{
			Command: "echo 'hello world'",
			Stdout:  stdout,
		})
		assert.NoError(t, err)

		assert.Equal(t, "hello world\n", stdout.String())
	})

	t.Run("can read stderr", func(t *testing.T) {
		client, err := sshexec.NewClient(resource.GetHostPort("22/tcp"), "sshuser", sshkey)
		assert.NoError(t, err)

		stderr := &bytes.Buffer{}
		err = client.Run(&sshexec.Cmd{
			Command: "echo 'test' >&2",
			Stderr:  stderr,
		})
		assert.NoError(t, err)

		assert.Equal(t, "test\n", stderr.String())
	})
}
