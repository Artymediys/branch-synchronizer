# BSync – Branch Synchronizer

Приложение для синхронизации веток в проектах GitLab

## Конфигурационный файл
Пример конфига лежит в `./env/config.json`, а также представлен ниже:

```json
{
  "gitlab": {
    "url": "https://gitlab.example.com",
    "token": "your_gitlab_token",
    "group_id": 1234
  },
  "telegram": {
    "bot_token": "your_telegram_bot_token",
    "channel_id": -1001234567890
  },
  "mattermost": {
    "url": "https://mattermost.example.team",
    "bot_token": "your_mattermost_bot_token",
    "channel_id": "abde12fghijklmno3456pqr6s"
  },
  "branch_pairs": [
    "production -> rc",
    "rc -> master"
  ],
  "whitelist_projects": [
    "project1",
    "project2"
  ],
  "blacklist_projects": [],
  "check_interval_in_hours": 0
}
```

### Обязательные поля
Для работы приложения необходимо задать `gitlab` и `branch_pairs`:

```
gitlab:
url – Адрес **GitLab**
token – PAT (Personal Access Token) юзера **GitLab**
group_id – ID группы, в которой будут проверяться ветки проектов
```

```
branch_pairs – Пары веток, для которых будет происходить сравнение и создание Merge Requests
```
// _Ветки указываются в формате `sourceBranch -> targetBranch`.
Если ветки отличаются, то создастся **Merge Request** из `sourceBranch` в `targetBranch`_

### Необязательные поля
Эти поля необязательны к заполнению – `telegram`, `mattermost`, `whitelist_projects`, `blacklist_projects`, `check_interval_in_hours`:
```
telegram:
bot_token – Токен бота, от лица которого будут отправляться сообщения
channel_id – ID канала, куда будут отправляться сообщения
```
// _Если не требуется использование Telegram, то можно оставить поле пустим -> telegram: {}_

```
mattermost:
url – Адрес Mattermost
bot_token – Токен бота, от лица которого будут отправляться сообщения
channel_id – ID канала, куда будут отправляться сообщения
```
// _Если не требуется использование Mattermost, то можно оставить поле пустим -> mattermost: {}_


```
whitelist_projects – Список проектов (репозиториев), которые будут проверяться приложением
blacklist_projects – Список проектов (репозиториев), которые не будут проверяться приложением. То есть будут проверены все проекты в группе за исключением этого списка
```
// _Названия проектов указываются в формате "url-названия". 
То есть если url проекта выглядит так: https://gitlab.example.com/some-group/some-project, 
то названием здесь является последняя часть url => `some-project`_  
// _Если не указать ни одно из этих двух полей, то приложение будет проверять все проекты в группе_  
// _Если одновременно в двух полях будут указаны значения, то приложение завершится с ошибкой_


```
check_interval_in_hours – Интервал проверок, указывающийся в часах.
```
// _Если не указать это поле или указать значение `<= 0`, то приложение выполнит проверку 1 раз_

## Запуск
Приложение можно запустить, используя заранее собранные исполняемые файлы, находящиеся в директории `./bin`.

Также есть вариант собрать приложение самостоятельно, как это описано в разделе [Сборка](#сборка).

Для запуска приложения, требуется указать конфигурационный файл с помощью флага `-config`.  

**Пример:**
```shell
./bin/linux/bsync_amd64 -config ./env/config.json
```

После запуска приложения, оно генерирует лог-файл, который будет располагаться по пути `./log/bsync.log`
относительно исполняемого файла.

## Сборка
Требуется установленный **Go** версии `>=1.22.5`.

В корне репозитория находится bash-скрипт `build.sh` для сборки приложения под необходимую ОС и архитектуру.

Скрипт можно запустить вместе с аргументами:
1) Сборка приложения под предустановленные ОС и архитектуры, а именно –> **windows amd64**; **linux amd64**;
   **macos amd64** и **arm64**
```shell
./build.sh all
```

2) Сборка приложения под конкретную ОС. Из доступных вариантов -> **windows**; **linux**, **macos**
```shell
./build.sh <ОС>
```

Также скрипт можно запустить без аргументов. Тогда он соберёт приложение под текущую ОС и архитектуру,
на которых был запущен скрипт:
```shell
./build.sh
```

Если требуется собрать приложение под ОС или архитектуру, которых нет в стандартном наборе скрипта `build.sh`,
тогда можно посмотреть список доступных для сборки ОС и архитектуры командой:
```shell
go tool dist list
```
после чего выполнить следующую команду для сборки:

**Linux/MacOS**
```shell
env GOOS=<ОС> GOARCH=<архитектура> go build -o <название итогового файла>
```

**Windows/Powershell**
```shell
$env:GOOS="<ОС>"; $env:GOARCH="<архитектура>"; go build -o <название итогового файла>
```

**Windows/cmd**
```shell
set GOOS=<ОС> && set GOARCH=<архитектура> && go build -o <название итогового файла>
```