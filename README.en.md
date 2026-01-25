# lai-panel

An application management panel built with Go and Vue 3, supporting Docker container management, node management, and application deployment.

## Features

- 🚀 **Application Management** - Create, edit, and manage applications
- 🖥️ **Node Management** - Manage local and remote server nodes
- 🐳 **Docker Management** - Container, image, network, and volume management
- 📦 **Service Deployment** - Deploy services via Docker Compose
- 💻 **Terminal Access** - Access nodes and containers through web terminal
- 📁 **Workspace Management** - Manage application workspace files
- 🔐 **User Authentication** - JWT-based user authentication system

## Tech Stack

### Backend
- Go 1.25+
- CloudWeGo Hertz (HTTP framework)
- SQLite (Database)
- Docker API
- SignalR (Real-time communication)

### Frontend
- Vue 3
- PrimeVue 4
- TypeScript
- Vite
- TailwindCSS

## Quick Start

### Prerequisites

- Go 1.25 or higher
- Node.js 18+ and npm
- Docker (optional, for container management)

### Installation

1. Clone the repository
```bash
git clone https://github.com/benlocal/lai-panel.git
cd lai-panel
```

2. Build the project
```bash
make all
```

3. Build the frontend
```bash
cd dashboard
npm install
npm run build
cd ..
```

4. Run the service
```bash
make run-serve
```

The service runs on `http://localhost:8080` by default.

### Development Mode

Run the frontend development server:
```bash
cd dashboard
npm run dev
```

The frontend development server runs on `http://localhost:5173`

## Project Structure

```
lai-panel/
├── cmd/
│   ├── serve/      # Main service program
│   └── agent/       # Agent program
├── dashboard/       # Frontend Vue application
├── pkg/            # Go packages
│   ├── handler/    # HTTP handlers
│   ├── repository/ # Data access layer
│   ├── model/      # Data models
│   └── ...
└── migrations/     # Database migration files
```

## Usage

1. Access `http://localhost:8080` to open the management panel
2. Login with default credentials (first run requires user creation)
3. Add nodes, create applications, and start deploying services

## License

See the [LICENSE](LICENSE) file for details.
