#!/bin/bash
git pull --ff-only
make build
cd ./dist
./recipya migrate up
sudo supervisor stop recipya
sudo supervisor start recipya
