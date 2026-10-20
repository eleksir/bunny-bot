# Bunny telegram bot

Ещё один телеграм-бот.

Его задача - банить спаммеров. Обнаружение спаммеров работает единственным способом - если вошедший не отвечает на
приветствие в течение 15 минут, он попадает в бан.

## Как это собрать?

Понадобится компилятор go версии 1.24 или новее. Воодим в каталоге с исходниками команду

```bash
make
```

и ждём. Через некоторое время в каталоге должен появиться бмнапник bunny-bot или bunny-bot.exe - это и есть весь наш бот.

## Как это запустить

Надо создать каталог `data/config.json` рядом с бинарником в качестве примера можно использовать `data/config-example.json`:

```json
{
  "token": "get_it_from_botfather",
  "loglevel": "debug",
  "logfile": "/var/log/bunny-bot/bunny-bot.log",
  "chat_logs": "chat_logs"
}
```

Затем запускаем сам бинарник:

```bash
./bunny-bot
```

Путь до конфига захардкожен и это `<binary dir>/data/config.json`.

### Логгирование чатов

Параметр `chat_logs` опционален и, если он есть и целевой каталог существует и доступен для записи, бот будет логгировать события в чате.

События из чата группируются по дням:

```text
<chat_logs>/<chat id>/<yyyy>/<MM>/<dd>.txt
```

И каждая строка выглядит примерно так:

```text
21:33:07.123 in chat -1001234567890 (title=Bunny Fans username=@bunnybot) user 42 (first_name=Ann last_name=Smith username=@ann) says hello there
21:33:07.123 in chat -1001234567890 (title=Bunny Fans username=@bunnybot) user 42 (first_name=Ann last_name=Smith username=@ann) post video (file_name=clip.mp4 mime=video/mp4 1280x720 duration=15s size=1048576) says look at this reply_to=1
21:33:07.123 in chat -1001234567890 (title=Bunny Fans username=@bunnybot) user 42 (first_name=Ann) user 7 (first_name=Al) joined the chat
21:33:07.123 in chat -1001234567890 (title=Bunny Fans) user 42 (first_name=Ann) status member -> administrator custom_title=Mod by user 7 (first_name=Bob)
```

Временные метки во всех логах проставляются по часовому поясу GMT+0.

### Дополнительные файлы

в каталоге `contrib` есть пример openrc-скрипта для запуска бота в среде alpine linux.

### Signals

| Signal | Effect |
|---|---|
| `SIGINT`, `SIGTERM`, `SIGQUIT` | Аккуратоное запершение работы - закрываются кэши, логи и затем программа завершается. |
| `SIGHUP`, `SIGUSR1` | Переоткрывает лог-файлы, для ротации логов. Процесс продолжает работать. |

N.B.
Бот не откладывает PID-файла, используйте copy-truncate в logrotate.

## Special thanks

To all active users of chat :)

To maintainers of telegram bot lib [https://github.com/NicoNex/echotron](https://github.com/NicoNex/echotron)

And kind hello to maintainer of spew dumper [https://github.com/davecgh/go-spew](https://github.com/davecgh/go-spew)

Without you, folks, all this bot-thing become very boring adventure.

## AI Slop warning

Последние изменения в боте сделаны с помощью ai.
