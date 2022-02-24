FROM swaggerapi/swagger-ui:latest
COPY ./openapi /usr/share/nginx/html/api
ENV URLS="[{url: 'api/merged.swagger.json', name: 'merged'},{url: 'api/user.swagger.json', name: 'user'},{url: 'api/saas.swagger.json', name: 'saas'},{url: 'api/sys.swagger.json', name: 'sys'}]"