#!/bin/bash
git pull --ff-only
make build-deploy
cd ./dist
./recipya migrate up
sudo supervisorctl stop recipya
sudo supervisorctl start recipya
