FROM node
WORKDIR /messaging
COPY package.json package-lock.json ./
RUN npm install
COPY . .
EXPOSE 5001-5003
ENTRYPOINT ["node", "index.js"]