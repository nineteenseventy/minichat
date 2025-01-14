FROM docker.io/nginx:stable-alpine

# this image exposes port 8080

COPY dist/web /usr/share/nginx/html/
COPY web/nginx.conf /etc/nginx/conf.d/default.conf
COPY web/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh


ENTRYPOINT [ "/entrypoint.sh" ]
CMD ["nginx", "-g", "daemon off;"]