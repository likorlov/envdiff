# envdiff

> CLI tool to diff and reconcile environment variable files across deployment stages

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git && cd envdiff && go build -o envdiff .
```

## Usage

Compare two `.env` files and see what's missing or changed between deployment stages:

```bash
# Diff two env files
envdiff diff .env.staging .env.production

# Show only keys missing in the target file
envdiff diff --missing .env.staging .env.production

# Reconcile by copying missing keys from source to target
envdiff reconcile .env.staging .env.production

# Output diff in JSON format
envdiff diff --format json .env.staging .env.production
```

### Example Output

```
Missing in .env.production:
  - NEW_FEATURE_FLAG
  - CACHE_TTL

Changed values:
  ~ DATABASE_URL  (staging != production)

Extra in .env.production:
  + LEGACY_API_KEY
```

## Why envdiff?

Managing environment variables across `development`, `staging`, and `production` is error-prone. `envdiff` makes it easy to spot configuration drift before it becomes a production incident.

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

[MIT](LICENSE)