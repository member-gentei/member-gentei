# add gentei-member-check system user
getent group gentei-member-check >/dev/null || groupadd -r gentei-member-check
getent passwd gentei-member-check >/dev/null || \
    useradd -r -g gentei-member-check -s /sbin/nologin \
    -c "gentei-member-check.service user" gentei-member-check
exit 0