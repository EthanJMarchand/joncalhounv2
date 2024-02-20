FROM node:latest AS tailwind-builder
WORKDIR /tailwind
RUN npm init -y && \
    npm install tailwindcss && \
    npx tailwindcss init
COPY ./templates /templates
COPY ./tailwind/tailwind.config.js /src/tailwind.config.js
COPY ./tailwind/styles.css /src/styles.css
RUN npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /styles.css --minify



# FROM is the image you want to build the conatiner from. AS simply gives the container a name
FROM golang:alpine AS builder
# WORKDIR gives the container a root directory name
WORKDIR /app
# COPY copies the first part and every part after, to the last argument directory on the container
COPY go.mod go.sum ./
# RUN is a command within the container. The string aftewards is like we are typing in the actual comand line
RUN go mod download
# COPY this line is basically taking everything ( . ) and copying it to the root directory (/app) in the container.
COPY . .
# RUN is a command within the container. -0 (Stands for output). It is the server binary in the current directory of the container.
# ./cmd/server/ is where we are building from on our local machine.
RUN go build -v -o ./server ./cmd/server/



FROM alpine
WORKDIR /app
COPY ./assets ./assets
COPY .env .env
COPY --from=builder /app/server ./server
COPY --from=tailwind-builder /styles.css app/assets/styles.css
# CMD actually runs binary files.
CMD ./server