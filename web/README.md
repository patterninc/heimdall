# Heimdall Web UI

## Local Development Setup

### Prerequisites

- Node.js
- PNPM (Package manager)

### Installing Dependencies

```bash
cd web
pnpm install
```

### Running the Application

1. First, ensure cookie-monster is running (see Cookie Monster setup below)
2. Start the development server:

```bash
pnpm dev
```

The application will be available at http://localhost:4000

### Setting up Cookie Monster

Cookie Monster is required for local development to handle authentication and API proxying. Follow these steps to set it up:

1. Download cookie-monster from [Gatekeeper Releases](https://github.com/patterninc/gatekeeper/releases/tag/v1.0.3)
2. Unzip the downloaded file
3. Run cookie-monster using the following command:
   ```bash
   cookie-monster-darwin-arm64 -relay 9090
   ```

#### First Time Setup Notes:

- Since the executable is not signed yet, you'll need to:
  1. Go to System Settings > Security & Privacy on your Mac
  2. Allow the cookie-monster executable to run

#### What happens when you run cookie-monster:

1. First Run:

   - It will install Chromium browser
   - Open a window with Microsoft Entra ID SSO prompt
   - You'll need to:
     - Provide your credentials and login
     - Complete 2FA challenge
   - After successful login, it creates a cookie file at:
     - `$HOME/.pattern/gatekeeper/heimdall.json`

2. Subsequent Runs:
   - If cookies are still valid, it will just start the relay
   - The relay proxy will point to https://heimdall.aws.pattern.com
   - Your local web UI can now connect to http://127.0.0.1:9090
