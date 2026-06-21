# Linux Only

## Why Linux Only?

The first version of `gofplll` (v0.1) supports only Linux for these reasons:

1. **fplll is primarily developed for Linux.** Pre-built binaries are readily available on Debian, Ubuntu, Arch, and other Linux distributions.

2. **WSL2 support.** Most Go developers on Windows use WSL2, which provides a full Linux environment. This library works seamlessly under WSL2.

3. **Simpler subprocess model.** Linux `exec.CommandContext` provides reliable process termination via `SIGKILL` on timeout. Cross-platform process management adds complexity that is unnecessary for v0.1.

4. **No cgo dependency.** By wrapping the `fplll` binary rather than linking `libfplll`, we avoid cgo entirely. This makes cross-compilation straightforward and eliminates platform-specific build issues.

## Building for Non-Linux

On non-Linux systems, the `fplll_unsupported.go` file is compiled. It provides the same API surface but returns `ErrUnsupportedOS` for all operations. This allows your code to compile on any platform, but runtime calls will fail gracefully.

```go
client, err := gofplll.NewDefault()
// err == gofplll.ErrUnsupportedOS on non-Linux
```

## Future Plans

A future version may support macOS and other Unix-like systems where fplll is available.
