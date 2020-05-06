FROM node:14.0.0

WORKDIR /app
COPY package.json /app
RUN npm install && mkdir /app/storage
COPY . /app

CMD node index.js
