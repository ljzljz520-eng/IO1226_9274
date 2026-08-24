# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	miniarrow/cmd/archer	[no test files]
ok  	miniarrow/internal/difficulty	0.002s
ok  	miniarrow/internal/engine	0.001s
ok  	miniarrow/internal/model	0.001s
ok  	miniarrow/internal/report	0.001s
ok  	miniarrow/internal/scheduler	0.002s
--- FAIL: TestWorkflow27 (0.02s)
    regression_test.go:30: summary disappeared after consecutive dispatch: first=1 second=0
FAIL
FAIL	miniarrow/internal/service	0.028s
ok  	miniarrow/internal/store	0.014s
ok  	miniarrow/internal/upgrade	0.001s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/archer): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/archer): exit `0`
- Frontend build (web): exit `0`
