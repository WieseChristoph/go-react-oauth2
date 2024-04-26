# Go React OAuth2

[![Go Report Card](https://goreportcard.com/badge/github.com/WieseChristoph/go-react-oauth2/backend)](https://goreportcard.com/report/github.com/WieseChristoph/go-react-oauth2/backend)

This is a simple full stack application that uses a Go backend with OAuth2 to authenticate users with Discord and Google. The frontend is a React application built with Vite. The whole application with SSL certificate generation is dockerized.

It is based on [go-oauth2-backend](https://github.com/WieseChristoph/go-oauth2-backend).

## Environment Variables

Copy the `.env.example` file to `.env` and fill in the values or set the environment variables manually.

### General
- `APP_URL` - Application URL (with port if needed) (used for redirecting)

### Database
- `DB_DATABASE` - Database name
- `DB_USERNAME` - Database username
- `DB_PASSWORD` - Database password

### OAuth2
- `DISCORD_CLIENT_ID` - Discord OAuth2 client ID
- `DISCORD_CLIENT_SECRET` - Discord OAuth2 client secret
- `GOOGLE_CLIENT_ID` - Google OAuth2 client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth2 client secret

### SSL
- `CERTBOT_EMAIL` - Email for Let's Encrypt certificate renewal
- `CERTBOT_DOMAIN` - Domain for Let's Encrypt certificate

## Domain renewal

To renew the SSL certificate automatically, add the following cron job to the server:

```
0 5  * * *  /path/to/project/.docker/renew_certificates.sh
```

## Development

To start the development environment with hot reloading, run the following command after setting the environment variables and installing the dependencies:

```bash
docker compose -f docker-compose-dev.yaml up
```

The application will be available at `http://localhost:8080`.

## Production

To start the production environment, run the following command:

```bash
docker compose up
```

If the DNS records are correctly set up, the application will be available at the domain specified in the `.env` file.
