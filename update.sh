#!/bin/bash
git pull --ff-only
make build-deploy
cd ./bin
./recipya migrate up
sudo supervisorctl stop recipya
sudo supervisorctl start recipya
