FROM node:latest

WORKDIR /opt/agent
COPY index.js package.json ./

RUN npm install

CMD [ "node", "index.js" ]
