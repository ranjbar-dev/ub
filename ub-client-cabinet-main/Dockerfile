FROM node:18-slim
ENV NODE_OPTIONS=--openssl-legacy-provider
WORKDIR /usr/src/app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile --ignore-optional
COPY . ./
RUN yarn build-dev
