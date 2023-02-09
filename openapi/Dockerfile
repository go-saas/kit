FROM swaggerapi/swagger-ui:latest
COPY ./openapi /usr/share/nginx/html/api
ENV URLS="[{url: 'api/kit-merged.swagger.json', name: 'kit-merged'},{url: '/v1/user/dev/swagger/swagger_spec', name: 'user'},{url: '/v1/saas/dev/swagger/swagger_spec', name: 'saas'},{url: '/v1/sys/dev/swagger/swagger_spec', name: 'sys'},{url: '/v1/realtime/dev/swagger/swagger_spec', name: 'realtime'}]"
EXPOSE 8080