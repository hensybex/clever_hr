# Stage 1: Build the Flutter web app
FROM ghcr.io/cirruslabs/flutter:stable AS build

WORKDIR /app

# Copy the pubspec files and get dependencies
COPY pubspec.* ./
RUN flutter pub get

# Copy the rest of the app source code
COPY . .

# Build the web app without analyzer
RUN flutter pub get; flutter build web --release

# # Stage 2: Serve the app with nginx
FROM nginx:stable-alpine

# # Copy the build output to nginx's html directory
COPY --from=build /app/build/web /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

# # Expose port 80
EXPOSE 80

# # Start nginx
CMD ["nginx", "-g", "daemon off;"]
