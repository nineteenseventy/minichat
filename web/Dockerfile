FROM docker.io/nginxinc/nginx-unprivileged:stable-alpine

# this image exposes port 8080


COPY dist/web /usr/share/nginx/html/
COPY web/nginx.conf /etc/nginx/conf.d/default.conf
CMD ["nginx", "-g", "daemon off;"]