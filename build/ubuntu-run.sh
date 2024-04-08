#!/bin/bash

export FQDN

start_xrdp_services() {
    # Preventing xrdp startup failure
    rm -rf /var/run/xrdp-sesman.pid
    rm -rf /var/run/xrdp.pid
    rm -rf /var/run/xrdp/xrdp-sesman.pid
    rm -rf /var/run/xrdp/xrdp.pid
	rm -rf /var/run/sssd.pid
	service dbus restart
	service xrdp restart
    # Use exec ... to forward SIGNAL to child processes
    #xrdp-sesman && xrdp
}

stop_xrdp_services() {
    pkill xrdp
	pkill dbus
}

echo "Script do entryponit está executando..."

#while :; do
#	authUp=$(curl -skS "https://$FQDN/auth/realms/master/protocol/saml/descriptor" 2>&-)
#	[[ $authUp =~ "entityID" ]] && break || ( echo -e "Aguardando o Auth Jump..."; sleep 5 )
#done

    addgroup $workUser
    
    useradd -m -s /bin/bash -g $workUser $workUser
    wait

    echo $workUser:$workPass | chpasswd 
    wait

    if [[ $root == "yes" ]]; then
        usermod -aG sudo $workUser
    fi
    wait
    echo "User '$workUser' is added"

echo -e "Usuário padrão criado...\n"

trap "stop_xrdp_services" SIGKILL SIGTERM SIGHUP SIGINT EXIT

#/etc/config.sh

echo -e "Inicializando o XRDP...\n"

#[ -z "${DTL+x}" ] && DTL=0
#[ -z "${ITL+x}" ] && ITL=0

#sed -i "s/KillDisconnected=false/KillDisconnected=true/;s/DisconnectedTimeLimit=0/DisconnectedTimeLimit=$DTL/;s/IdleTimeLimit=0/IdleTimeLimit=$ITL/" /etc/xrdp/sesman.ini

stop_xrdp_services
start_xrdp_services

restartXrdp > /var/log/restartXrdp.log 2>&1 & tail -f /dev/null