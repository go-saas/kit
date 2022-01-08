FROM swaggerapi/swagger-ui:latest
COPY ./openapi /usr/share/nginx/html/api
ENV URLS="[{url: 'api/user.swagger.json', name: 'user'},{url: 'api/saas.swagger.json', name: 'saas'}]"