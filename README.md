# sendmail - GO Mail Service

A lightweight HTTP service to send emails via any SMTP server (e.g. Gmail, Outlook, Postfix).
Designed to be minimal and easy to deploy with Docker.

It exposes a single API endpoint with a JSON payload:

```
POST /v1/send
```

```json
{
  "from": "your@gmail.com",
  "to": ["recipient@example.com"],
  "subject": "Test email",
  "body": "Hello sendmail"
}
```

Authentication is handled via a **Bearer token**.

---

## Installation

### 1. Go
Clone the repo and build:

```bash
git clone https://github.com/inis17/sendmail-docker.git
cd sendmail-docker.git
go build -o sendmail main.go
```

Run the server:

```bash
SMTP_HOST=smtp.gmail.com \
SMTP_PORT=587 \
SMTP_USER=your@gmail.com \
SMTP_PASS=your-app-password \
API_TOKEN=mysecrettoken \
./sendmail
```


---

### 2. Docker

Run:

```bash
docker run --rm -p 8080:8080 \
  -e SMTP_HOST=smtp.gmail.com \
  -e SMTP_PORT=587 \
  -e SMTP_USER=your@gmail.com \
  -e SMTP_PASS=your-app-password \
  -e API_TOKEN=mysecrettoken \
  inis17/sendmail
```

---

### 3. Docker Compose

Example `docker-compose.yml`:

```yaml
services:
  sendmail:
    image: inis17/sendmail
    ports:
      - "8080:8080"
    environment:
      - SMTP_HOST=smtp.gmail.com
      - SMTP_PORT=587
      - SMTP_USER=your@gmail.com
        # tou can append _FILE to load from a file inside the container (e.g. docker secrets)
      - SMTP_PASS=yourapppassword
      - API_TOKEN=longrandomstring
```

---

## Usage

```bash
curl -X POST http://localhost:8080/v1/send \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer mysecrettoken" \
  -d '{
        "from": "your@gmail.com",
        "to": ["your@gmail.com"],
        "subject": "Homelab Alert",
        "body": "Something went wrong!"
      }'
```

---

## Environment Variables

| Variable         | Default  | Description                                                                   |
| ---------------- | -------- | ----------------------------------------------------------------------------- |
| `SMTP_HOST`      | *(none)* | SMTP server hostname (e.g. `smtp.gmail.com`)                                  |
| `SMTP_PORT`      | `587`    | SMTP server port (`587` for STARTTLS, `465` for implicit TLS, `25` for plain) |
| `SMTP_USER`      | *(none)* | SMTP username (usually your email address)                                    |
| `SMTP_PASS`      | *(none)* | SMTP password or app password                                                 |
| `SMTP_PASS_FILE` | *(none)* | Path to file containing SMTP password (overrides `SMTP_PASS`)                 |
| `API_TOKEN`      | *(none)* | Bearer token for authentication                                               |
| `API_TOKEN_FILE` | *(none)* | Path to file containing API token (overrides `API_TOKEN`)                     |
| `PORT`           | `8080`   | HTTP listen port                                                              |

---

## Security Notes

* Use **App Passwords** when connecting to Gmail/Outlook.
* Always use HTTPS (behind a reverse proxy like Nginx or Caddy) if exposed outside localhost.
* Keep your `API_TOKEN` long and random.

---

## Features

* Simple `/v1/send` endpoint
* Authentication via Bearer token
* Works with Gmail, Outlook, or any SMTP server
* Supports Docker secrets (`*_FILE`)
* Tiny, single binary service (\~10 MB with Go)
