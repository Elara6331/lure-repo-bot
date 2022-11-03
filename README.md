# LURE Repo Bot

A Github bot that reviews PRs to the LURE repo by analyzing the script for errors and providing comments on how to fix them.

There is also a command-line tool at `./cmd/lure-analyzer` that does the same thing but as a command.

## Configuration

### `LURE_BOT_ADDR`

The listen address for the webhook server. `:8080` by default.

### `LURE_BOT_GITHUB_TOKEN`

The Github token to be used for writing PR reviews

### `LURE_BOT_SECRET`

The secret used when setting up the Github webhook, used to verify the authenticity of webhook data.