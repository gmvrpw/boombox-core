# Boombox - modular, open and easily customizable music player designed with respect for the ideas of clean architecture.

> Boombox Core **is not a standalone piece of code** and it does not work by itself.
> This repository contains only player features (working with requests and users)
and does not know how to play and broadcast audio.
> You can find a fully working examples in the [corresponding repository](https://github.com/gmvrpw/boombox).

## Environment Variables
Environment variables are used only for passing secrets, and the path to the configuration file, but not the program configuration itself.

All secrets have a paired environment variable with the suffix `_FILE`,
which specifies the path to the file containing the secret.
The secret file is less specific than the environment variable
and will not be used if the secret's environment variable exists.

| Secret                               | Description                                |
|--------------------------------------|--------------------------------------------|
| `BOOMBOX_CONFIG_FILE`                | Config file path                           |
| `BOOMBOX_DISCORD_BOT_TOKEN`          | Discord front-end api token                |
| `BOOMBOX_REQUESTS_POSTGRES_USER`     | User for requests storage (PostgreSQL)     |
| `BOOMBOX_REQUESTS_POSTGRES_PASSWORD` | Password for requests storage (PostgreSQL) |
| `BOOMBOX_REQUESTS_POSTGRES_DB`       | Database for requests storage (PostgreSQL) |

## Configuration
The parameters specified in the config file are the most specific and if they are present, the program will not use environment variables.
*- Yes, this means that you can specify secrets in the config file, but we strongly discourage you from doing so.*

### Internal

| Field              | Type     | Description                                                                       |
|--------------------|----------|-----------------------------------------------------------------------------------|
| `requests.url`     | `string` | PostgreSQL DSN. Format: `postgresql://[user[:password]@][netloc][:port][/dbname]` |

### Runners

In the runners part, you can specify an array of entities configured for use in the run.

*Unfortunately, it is not currently possible to add runners in the runtime, although this feature is present in our roadmap.
That is why you should configure here all the runners you want to use.*

| Field              | Type     | Description                                                                                                         |
|--------------------|----------|---------------------------------------------------------------------------------------------------------------------|
| `runners[n].name`  | `string` | Runner's name                                                                                                       |
| `runners[n].url`   | `string` | URL to the runner's API                                                                                             |
| `runners[n].owner` | `string` | (Optional) Runner's owner.                                                                                          |
| `runners[n].test`  | `string` | (Optional) A regular expression that the request URL should match before checking playability on the runner's side. |
