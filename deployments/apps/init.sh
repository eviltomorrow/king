#!/bin/bash

mkdir -p $(pwd)/data/{king-collector,king-email,king-repository}/log
mkdir -p $(pwd)/data/king-email/etc

# smtp.conf
cat > $(pwd)/data/king-email/etc/smtp.json <<EOF
{
    "server":"mail.liarsa.me",
    "port":587,
    "username":"assistant@liarsa.me",
    "password":"5r6WAmzs2xyGMPqB",
    "alias":"assistant"
}
EOF
