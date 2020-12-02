# add gentei-bot system user
getent group gentei-bot >/dev/null || groupadd -r gentei-bot
getent passwd gentei-bot >/dev/null || \
    useradd -r -g gentei-bot -s /sbin/nologin \
    -c "gentei-bot.service user" gentei-bot
exit 0