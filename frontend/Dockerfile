# Étape 1 : Build Flutter web
FROM ghcr.io/cirruslabs/flutter:3.32.5 AS build

WORKDIR /app
COPY pubspec.yaml pubspec.lock ./
RUN flutter pub get
COPY . .
RUN flutter build web

# Étape 2 : Nginx pour servir le site
FROM nginx:alpine
COPY --from=build /app/build/web /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
