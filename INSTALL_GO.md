# Install Go

## Option 1: Install via Homebrew (Recommended for macOS)

```bash
brew install go
```

After installation, verify:
```bash
go version
```

You should see something like: `go version go1.21.x darwin/arm64`

## Option 2: Download from Official Website

1. Visit: https://go.dev/dl/
2. Download the macOS installer (.pkg file)
3. Run the installer
4. Verify installation: `go version`

## After Installation

Once Go is installed, you can run the backend:

```bash
cd backend
go mod download
go run cmd/main.go --api --port=8080
```

## Verify Go Installation

```bash
# Check Go version
go version

# Check Go environment
go env

# Test Go installation
go run --help
```

## Troubleshooting

**If `go` command is not found after installation:**

1. Add Go to your PATH. Add this to your `~/.zshrc`:
   ```bash
   export PATH=$PATH:/usr/local/go/bin
   # Or for Homebrew installation:
   export PATH=$PATH:$(brew --prefix)/opt/go/libexec/bin
   ```

2. Reload your shell:
   ```bash
   source ~/.zshrc
   ```

3. Verify:
   ```bash
   which go
   go version
   ```

**For Homebrew installation**, the PATH is usually set automatically, but if not, check:
```bash
echo $PATH | grep go
```

If Go is installed but not in PATH, add:
```bash
export PATH=$PATH:$(go env GOROOT)/bin
```

