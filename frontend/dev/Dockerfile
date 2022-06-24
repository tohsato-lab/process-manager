### ベースステージ ###
FROM node:16.15.1 as base

# install angular
RUN npm install -g @angular/cli

WORKDIR /app

COPY ./app/package*.json /app/

RUN rm -rf node_modules && npm install

### ビルドステージ ###
FROM base as dev-build

# 全てのソースファイルをコピー
COPY ./app /app/
