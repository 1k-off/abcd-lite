# abcd-lite

**abcd - make deployment simple as an alphabet.**

A minimal, easy-to-use deployment service for IIS websites. abcd-lite is designed for simplicity, focusing on essential features and minimal configuration. Perfect for teams or individuals who want to deploy and manage IIS sites with as little friction as possible.

---

# TL; DR
- Download and unpack latest release
- Generate config (`./abcd-lite.exe config generate`) and store generated password.
- Add IIS reverse-proxy configuration. Add `app.allowed_origins` to config file.
- Start it (`./abcd-lite.exe run`). Install windows service if needed using NSSM.
- Use GitHub action, Azure DevOps action, or whatever action you need.

---

## Features
- Deploy IIS websites
- Simple, modern web UI
- Single binary deployment
- Authenticated admin access
- API key management for projects

---

## Quick Start

### Development

1. **Start the backend:**
   ```powershell
   go run . start
   ```
2. **Start the frontend (in a separate terminal):**
   ```powershell
   cd frontend
   pnpm install
   pnpm dev
   ```
   The frontend will run on its own port (e.g., http://localhost:5173) and proxy API requests to the backend.

### Production (Single Binary with Embedded Frontend)

1. **Build the frontend:**
   ```powershell
   cd frontend
   pnpm build
   cd ..
   ```
2. **Build the Go binary with embedded frontend:**
   ```powershell
   go build -tags=embed_frontend
   ```
3. **Run the binary:**
   ```powershell
   ./abcd-lite.exe
   ```
   The app will serve both the API and the frontend from a single process.

---

## Usage
- Open the web UI in your browser.
- Log in with your admin token.
- Create and manage IIS website projects.
- Generate and manage API keys for deployments.
- Deploy your site with minimal hassle.

---

## Configuration
- All settings are in a single config file (see `configs/` or environment variables).
- Only the essentials: admin token, JWT secret, allowed origins, and database path.

---

## License
MIT + Country Restriction