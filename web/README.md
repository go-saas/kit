## Install and use

- Get the project code `

- Installation dependencies

```bash

pnpm install

```

- run

```bash
yarn serve
```

- build

```bash
yarn build
```

### API generator

```
openapi-generator-cli generate -i ../openapi/merged.swagger.json -g typescript-axios -o ./src/api-gen/ -c ./tools/config.json

```
