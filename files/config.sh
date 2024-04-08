#!/bin/bash

sed 's/$vaultIP/'"$FQDN"'/g' /opt/firefox/distribution/policies.json.model > /opt/firefox/distribution/policies.json

echo -e "FQDN="$FQDN"
auth_admin_pass="$auth_admin_pass"
auth_admin_user="$auth_admin_user"" > /tmp/var.env

cat /tmp/var.env >> /etc/environment

cp /tmp/pam_kc.so /lib/x86_64-linux-gnu/security/pam_kc.so
cp /tmp/libnss_ldh.so.2	/lib/x86_64-linux-gnu/libnss_ldh.so.2