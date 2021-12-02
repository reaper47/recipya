#!/bin/bash
git pull --ff-only
make build-deploy
cd ./dist
./recipya migrate up
sudo supervisor stop recipya
sudo supervisor start recipya
