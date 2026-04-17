# E2E (Playwright headed)

## Install

```bash
cd e2e
npm install
npx playwright install chromium
```

## Run (visible browser)

Start the API server separately (defaults to `http://localhost:5000`), then:

```bash
cd e2e
npm run test:headed
```

### Options

- Use a different server URL:

```bash
cd e2e
PLAYWRIGHT_BASE_URL="http://localhost:5000" npm run test:headed
```

- Make actions slower/easier to see:

```bash
cd e2e
PLAYWRIGHT_SLOWMO_MS=200 npm run test:headed:slow
```

