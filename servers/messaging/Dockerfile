
# use `node` image as the base
FROM node
# set the current working directory to /app
WORKDIR /app
# copy the package.json and package-lock.json files to the work directory
COPY package.json package-lock.json ./
# run npm install to install all dependent packages
RUN npm install
# copy your JavaScript source files
COPY . .
# set the entrypoint command
ENTRYPOINT ["node", "index.js"]