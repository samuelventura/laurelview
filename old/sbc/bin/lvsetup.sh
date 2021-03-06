#!/bin/bash -xe

chmod a+x bin/*
chmod 644 .ssh/*
chmod 600 .ssh/id_rsa

echo "LV_NBE_ENDPOINT=0.0.0.0:80" > /usr/local/bin/lvnss.config
echo "LV_NUP_ENDPOINT=127.0.0.1:80" >> /usr/local/bin/lvnss.config
echo "LV_DPM_ENDPOINT=127.0.0.1:81" >> /usr/local/bin/lvnss.config
echo "LV_NBE_DEBUG=127.0.0.1:82" >> /usr/local/bin/lvnss.config
echo "LV_NBE_DATABASE=/mnt/usb1/lvdb.db3" >> /usr/local/bin/lvnss.config
echo "LV_NSS_LOGS=/tmp" >> /usr/local/bin/lvnss.config
