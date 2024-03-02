FROM node:latest
WORKDIR /tailwind
RUN npm init -y && \
    npm install tailwindcss && \
    npm install -D @tailwindcss/forms && \
    npx tailwindcss init

CMD npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /dst/styles.css --watch --poll