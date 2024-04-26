# Go OAuth2 Backend

This is a simple backend written in Go that uses OAuth2 to authenticate users with Discord and Google. It can be extended to support other OAuth2 providers.

## Environment Variables

Copy the `.env.example` file to `.env` and fill in the values or set the environment variables manually.

### General
- `PORT` - Port to run the application
- `APP_URL` - Application URL (with port if needed) (used for redirecting)

### Database
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_DATABASE` - Database name
- `DB_USERNAME` - Database username
- `DB_PASSWORD` - Database password

### OAuth2
- `DISCORD_CLIENT_ID` - Discord OAuth2 client ID
- `DISCORD_CLIENT_SECRET` - Discord OAuth2 client secret
- `GOOGLE_CLIENT_ID` - Google OAuth2 client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth2 client secret

## MakeFile

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

Live reload the application
```bash
make watch
```

Clean up binary from the last build
```bash
make clean
```