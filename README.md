# Bunny telegram bot

Yet another telegram bot. It has one purpose - restrict access of abusers to chat.

By now it uses [Combot Anti-Spam System](https://cas.chat/) to detect spammers and its own method of detection of bad
users.

## Build instructions

Bot is not done yet, actually it does not do anything useful.

```bash
make
```

## Run instructions

You probably don't need it in its current state. But if you want - all is simple. Visit data/config.json, edit for your
needs and run:

```bash
./bunny-bot
```

There is openrc init script in contrib directory, you may use it for creating service under alpine linux or
gentoo(untested) with openrc init.

That's it.

## Special thanks

To all active users of chat :)

To maintainers of telegram bot lib [https://github.com/NicoNex/echotron](https://github.com/NicoNex/echotron)

To maintainers of pebble db [https://github.com/cockroachdb/pebble](https://github.com/cockroachdb/pebble)

To maintainers of zap logging facility [https://github.com/uber-go/zap](https://github.com/uber-go/zap)

And kind hello to maintainer of spew dumper [https://github.com/davecgh/go-spew](https://github.com/davecgh/go-spew)

Without you, folks, all this bot-thing become very boring adventure.

## Maybe some more mumble-jumble

Yes, but next time. When bot will be ready for some real action.

For now I can say that bot does not support any of commands.
