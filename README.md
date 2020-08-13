# discord-bot
Just a regular discord bot written in go

## Deployment
Bot automatically deployed when there are any changes to the `master` branch.

## Docker
Deployment success is based on if the docker image can get built or not

## Secrets
The bot is using `dotenv` to store the discord api token

# Run bot with updates
Use this `docker-compose.yml` file to run the bot while also recieving updates.  

`docker-compose.yml`
```yml
version : '3'
services:
  watchtower:
    container_name: watchtower
    restart: always
    image: containrrr/watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: discord-bot --debug true --interval 60
  discord-bot:
    container_name: discord-bot
    restart: always
    image: mehrn00/discord-bot-docker:latest
    environment:
     - DISCORD_TOKEN=<api-token>
     - PREFIX_CHARACTER=<single-character-prefix>
```
