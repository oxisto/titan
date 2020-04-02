# Titan

![build](https://github.com/oxisto/titan/workflows/build/badge.svg) 
[![](https://godoc.org/github.com/oxisto/titan?status.svg)](https://pkg.go.dev/github.com/oxisto/titan)

A simple go-based application that will help you become an industry titan in EVE Online.

## Register an application

In order to run `titan` for yourself, you need to register an application on the EVE Online Developer Portal (https://developers.eveonline.com/applications). Choose an appropriate name, such as `mytitan`. Connection type must be set to *Authentication & API Access*. The application needs to have the following permissions:
* publicData 
* esi-skills.read_skills.v1 
* esi-corporations.read_corporation_membership.v1 
* esi-ui.open_window.v1 
* esi-wallet.read_corporation_wallets.v1 
* esi-assets.read_corporation_assets.v1 
* esi-corporations.read_blueprints.v1 
* esi-industry.read_corporation_jobs.v1

In the last step, choose an appropriate redirect URI. This needs to include your hostname plus the path `/auth/callback`, if you want to deploy it somewhere or something like `http://localhost:4300/auth/callback` for a local development deployment.

## Set up development environment

To set up a simple development development environment with REDIS and MongoDB, you can use the included `dev.yml` file.

```
docker-compose -f dev.yml up -d
```

Titan uses [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) for configuring itself. This configuration options are defined in [server.go](https://github.com/oxisto/titan/blob/ec3aa32beb35dc13eb69a73445994aa02bf2eae4/cmd/server/server.go#L37-L57) and can be used in the following way:

- Environmental variables. The prefix is `TITAN`, so the config option for the redis server becomes `TITAN_REDIS`
- Command line arguments, such as `--redis`
- A configuration file, stored in `config/config.yaml`
