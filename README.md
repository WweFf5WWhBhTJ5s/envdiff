# envdiff

> CLI tool to diff and sync `.env` files across environments with secret masking support.

---

## Installation

```bash
go install github.com/youruser/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/envdiff.git && cd envdiff && go build -o envdiff .
```

---

## Usage

Compare two `.env` files and highlight differences:

```bash
envdiff diff .env.development .env.production
```

Sync missing keys from one file to another:

```bash
envdiff sync .env.development .env.production
```

Mask secret values when displaying output:

```bash
envdiff diff .env.staging .env.production --mask-secrets
```

Export the synced result to a new file:

```bash
envdiff sync .env.development .env.production --export .env.production.new
```

### Example Output

```
+ DB_HOST=localhost        (only in .env.development)
~ API_URL                  (value differs)
- STRIPE_KEY=***masked***  (only in .env.production)
```

---

## Flags

| Flag             | Description                                      |
|------------------|--------------------------------------------------|
| `--mask-secrets` | Redact sensitive values in output                |
| `--only-missing` | Show only keys missing from target               |
| `--export`       | Write synced output to a specified file          |
| `--no-color`     | Disable colored output (useful for CI pipelines) |

---

## License

[MIT](LICENSE) © 2024 youruser
