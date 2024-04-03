# Stage 1: Build
FROM node:14-alpine AS builder

WORKDIR /app

COPY package*.json ./

# Install dependencies including
RUN npm ci

COPY . .

RUN npm run build

# Stage 2: Run
FROM node:14-alpine

WORKDIR /app

RUN npm install -g serve

COPY --from=builder /app/build ./build

EXPOSE 3000

CMD ["serve", "-s", "build"]
