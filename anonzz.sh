#!/bin/bash
#Install and update machine
apt-get update
apt-get install snapd build-essential curl jq certbot wget -y
snap install core
snap refresh core
v=$(curl -s https://go.dev/dl/?mode=json | jq -r '.[0].version')
wget https://go.dev/dl/"${v}".linux-amd64.tar.gz
tar -C /usr/local -xzf "${v}".linux-amd64.tar.gz
ln -sf /usr/local/go/bin/go /usr/bin/go
rm "${v}".linux-amd64.tar.gz
ln -sf /usr/local/go/bin/go /usr/bin/go
if [ $# -ne 1 ]
then
  echo "Could not find domain. Specify it via parameter e.g. bash setup-ssl.sh DOMAIN.com"
  exit 1
fi

domain=$1

echo "Ssl in process"

certbot certonly --expand --manual --register-unsafely-without-email --agree-tos \
  --domain "${domain}" \
  --domain "*.${domain}" \
  --preferred-challenges dns

# certbot certificates

if [ $? -ne 0 ]
then
  echo "SSl failed for domain $domain"
  exit 1
fi
DefaultSSLDir="/etc/letsencrypt/archive"

if [ ! -d "$DefaultSSLDir" ]; then
  echo "Cannot Find Default SSL Directory /etc/letsencrypt/archive"
  exit 1
fi

certFile=`find /etc/letsencrypt/archive -type f -printf '%T@ %p\n' | sort -n | grep "cert" | tail -1 | cut -f2- -d" "`

privkeyFile=`find /etc/letsencrypt/archive -type f -printf '%T@ %p\n' | sort -n | grep "privkey" | tail -1 | cut -f2- -d" "`

mkdir -p /root/evilginx2/config/crt/$domain
cp $certFile /root/evilginx2/config/crt/$domain/o365.crt
cp $privkeyFile /root/evilginx2/config/crt/$domain/o365.key

echo "Finished"
