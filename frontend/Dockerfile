# Base
FROM node:25.2.1-alpine3.22 AS base
WORKDIR /app
COPY package*.json ./

# Development
FROM base AS dev
RUN npm install
CMD ["npm", "run", "dev"]

# Testing
FROM base AS test
RUN npm ci
COPY . .

# Build
FROM test AS build
RUN npm run build

# Production
FROM nginx:1.29.4-alpine3.23 AS prod
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]