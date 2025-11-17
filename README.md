<div align="center">
    <img 
        align="center" 
        src="https://raw.githubusercontent.com/robinvandernoord/tunnellm/master/_static/tunnelLM.png" 
        alt="Classy Configuraptor"
        width="400px"
        />
    <h1 align="center">TunnelLM</h1>
</div>


Your local OpenAI-compatible endpoint, powered by OpenRouter.
A tiny dockerized proxy that makes JetBrains IDEs think they're
talking to localhost while you're actually routing through OpenRouter.

## Why?

JetBrains IDEs accept a custom OpenAI URL but not a custom API key. 
TunnelLM solves this: point your IDE to `localhost:11434` and it proxies to OpenRouter with your key.

## Quick Start

```bash
# Option A: inline key
docker run -e OPENROUTER_API_KEY=sk-or-v1-... -p 11434:11434 robinvandernoord/tunnellm:latest

# Option B: env file
echo "OPENROUTER_API_KEY=sk-or-v1-..." > .env
docker run --env-file .env -p 11434:11434 robinvandernoord/tunnellm:latest
```

Now configure your IDE to use `http://localhost:11434` as the OpenAI endpoint.

## Configuration

| Variable             | Default                     | Notes                                      |
|----------------------|-----------------------------|--------------------------------------------|
| `OPENROUTER_API_KEY` | *required*                  | Your OpenRouter API key                    |
| `PORT`               | `11434`                     | Local port to bind                         |
| `ENDPOINT`           | `https://openrouter.ai/api` | OpenRouter API URL                         |
| `VERBOSITY`          | `0`                         | `0`=errors only, `1`=request/response logs |

## Development

Building from source requires `make` and `go`.
These will **not** be installed during `make install`, which only installs optional developer dependencies (`git-cliff`)

```bash
git clone git@github.com:robinvandernoord/tunnellm.git && cd tunnellm

make install  # dev deps, optional
make run      # build & run
```

## Docker Compose

```yaml
services:
  tunnellm:
    image: robinvandernoord/tunnellm:latest
    env_file: .env
    ports:
      - "11434:11434"
```
