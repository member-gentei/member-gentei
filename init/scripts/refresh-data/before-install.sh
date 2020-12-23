# add gentei-refresh-data system user
getent group gentei-refresh-data >/dev/null || groupadd -r gentei-refresh-data
getent passwd gentei-refresh-data >/dev/null || \
    useradd -r -g gentei-refresh-data -s /sbin/nologin \
    -c "gentei-refresh-data.service user" gentei-refresh-data
exit 0