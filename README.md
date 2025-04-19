# Containix

A terminal UI for managing Docker containers.

## Project Structure

The project follows a standard Go project layout:

```
containix/
├── cmd/                  # Command-line applications
│   └── containix/        # Main application entry point
├── internal/             # Private application code
│   ├── app/              # Application initialization
│   ├── docker/           # Docker client and operations
│   └── ui/               # User interface
│       ├── components/   # Reusable UI components
│       └── views/        # Application views/screens
├── pkg/                  # Public libraries (for potential reuse)
└── main.go               # Backward compatibility wrapper
```

## Development

### Building

```bash
go build -o containix ./cmd/containix
```

### Running

```bash
./containix
```

## Features

- List all Docker containers
- Start, stop, and restart containers
- View container logs
- Real-time updates

## Keyboard Shortcuts

- `s`: Stop the selected container
- `t`: Start the selected container
- `x`: Restart the selected container
- `l`: View logs of the selected container
- `r`: Refresh the container list
- `q`: Quit the application