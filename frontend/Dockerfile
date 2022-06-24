### ベースステージ ###
FROM node:16.15.1 as base

# install angular
RUN npm install -g @angular/cli

WORKDIR /app

COPY ./app/package*.json /app/

RUN rm -rf node_modules && npm install

### ビルドステージ ###
FROM base as build

# 全てのソースファイルをコピー
COPY ./app /app/

# 本番用ビルド
RUN ng build --configuration production --output-path=./dist/build-by-docker

### プロダクションステージ ###
FROM nginx:latest as prod

# ビルドステージで生成されたファイルをnginxの公開用ディレクトリにコピー
COPY --from=build /app/dist/build-by-docker /usr/share/nginx/html

# nginx.confをコピー
COPY ./nginx/nginx.conf /etc/nginx/nginx.conf
