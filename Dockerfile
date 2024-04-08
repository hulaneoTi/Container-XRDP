FROM ubuntu:23.04

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get -y update 
RUN apt-get -y upgrade

RUN apt-get install -y \
    xfce4 \
    xfce4-clipman-plugin \
    xfce4-cpugraph-plugin \
    xfce4-netload-plugin \
    xfce4-screenshooter \
    xfce4-taskmanager \
    xfce4-terminal \
    xfce4-xkb-plugin 

RUN apt-get install -y \
    sudo \
    wget \
    xorgxrdp \
    xrdp && \
    apt remove -y light-locker xscreensaver && \
    apt autoremove -y && \
    rm -rf /var/cache/apt /var/lib/apt/lists

RUN mkdir -p /var/run/dbus && \
    cp /etc/X11/xrdp/xorg.conf /etc/X11 && \
    sed -i "s/console/anybody/g" /etc/X11/Xwrapper.config && \
    sed -i "s/xrdp\/xorg/xorg/g" /etc/xrdp/sesman.ini && \
    echo "xfce4-session" >> /etc/skel/.Xsession
	
RUN apt update && apt install -y tar \
	bzip2 \
	nano \
	openbox \
	curl \
	jq \
	tmux \
	htop \
	pamtester \
	libpam-aad \
	libnss-aad \
	openssh-server \
	aad-cli \
	dbus-x11 \
	notepadqq

COPY ./build/firefox-latest.tar.bz2 /tmp/firefox-latest.tar.bz2
RUN tar -xjf /tmp/firefox-latest.tar.bz2 -C /opt/ && \
	ln -s /opt/firefox/firefox /usr/local/bin/ && \
    rm /tmp/firefox-latest.tar.bz2
	
COPY ./files/distribution/* /opt/firefox/distribution/

COPY ./build/ubuntu-run.sh /usr/bin/
RUN mv /usr/bin/ubuntu-run.sh /usr/bin/run.sh
RUN chmod +x /usr/bin/run.sh

COPY ./build/startwm-xfce.sh /etc/xrdp/startwm.sh
RUN chmod +x /etc/xrdp/startwm.sh

COPY ./build/aad.conf /etc/

COPY ./build/restartXrdp /usr/bin
RUN chmod +x /usr/bin/restartXrdp

RUN sudo sed -i '13i\session required pam_mkhomedir.so' /etc/pam.d/common-session

RUN echo "ALL ALL=(root) NOPASSWD: /bin/su" | sudo tee -a /etc/sudoers

RUN echo "setxkbmap -layout br" | sudo tee -a /etc/X11/xinit/xinitrc
RUN sudo sed -i 's/XKBLAYOUT="us"/XKBLAYOUT="br"/' /etc/default/keyboard
RUN echo "BROWSER=/opt/firefox/firefox" >> /etc/environment

COPY ./xrdp24b.jpeg /usr/share/wallpaper/
COPY ./xfce4-desktop.xml /etc/skel/.config/xfce4/xfconf/xfce-perchannel-xml/

ENTRYPOINT ["/usr/bin/run.sh"]
CMD ["$workUser", "$workPass", "$root"]