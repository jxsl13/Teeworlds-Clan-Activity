# Teeworlds Clan Activity

Create a file called `.env` in this directory and paste the following config file in there.  

```env
TCA_DISCORD_BOT_TOKEN=NjI..........

# channel where to write the online members message
TCA_DISCORD_CHANNEL=7896.............

# clan to track
TCA_TEEWORLDS_CLAN=[ friends ]
# interval between fetching the new playerlist
TCA_REFRESH_INTERVAL=60s

# leave empty if you do not want to use this service
TCA_DISCORD_NOTIFICATION_CHANNEL=7188..........
```

Then compile the application with:  

```shell
go get -d
go build .
```

and run it on some vps that is connected to the internet 24/7
