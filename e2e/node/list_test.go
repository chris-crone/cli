package node

import (
	"os"
	"testing"

	"github.com/gotestyourself/gotestyourself/assert"
	"github.com/gotestyourself/gotestyourself/icmd"
)

func TestNodeLs(t *testing.T) {
	cmd := icmd.Command("docker", "node", "ls")
	env := os.Environ()

	// Default orchestrator.
	result := icmd.RunCmd(cmd)
	result.Assert(t, icmd.Success)

	// Kubernetes.
	cmd.Env = envForOrchestrator(env, "kubernetes")
	result = icmd.RunCmd(cmd)
	result.Assert(t, icmd.Expected{ExitCode: 1})
	assert.Equal(t, result.Stderr(), "docker node ls is only supported on a Docker cli with either all or swarm features enabled\n")

	// Swarm.
	cmd.Env = envForOrchestrator(env, "swarm")
	result = icmd.RunCmd(cmd)
	result.Assert(t, icmd.Success)

	// All.
	cmd.Env = envForOrchestrator(env, "all")
	result = icmd.RunCmd(cmd)
	result.Assert(t, icmd.Success)
}

func envForOrchestrator(env []string, orchestrator string) []string {
	res := make([]string, len(env))
	copy(res, env)
	return append(res, "DOCKER_ORCHESTRATOR="+orchestrator)
}
