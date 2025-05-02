# syntax=docker/dockerfile:1

FROM node:23-slim AS build_image

ARG env_file=env/.env.dev

WORKDIR /app

# Copy files needed for dependency install
COPY client/package.json .
COPY client/yarn.lock .

# Install dependencies
RUN yarn install --immutable --immutable-cache --check-cache --omit=dev

# Copy the env file
COPY ${env_file} env/.env

# Copy the rest of the app
COPY client client

# Build the app
RUN cd client && yarn run build

FROM node:23-slim AS final_image

ARG port
ENV ENVPORT=${port}

# RUN apk add --no-cache ca-certificates

COPY --from=build_image /app/client/dist ./dist

RUN yarn global add serve

EXPOSE ${port}

CMD ["sh", "-c", "serve dist/ -l $ENVPORT"]