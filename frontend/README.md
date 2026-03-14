# Mission Control

### Run Dev
```
docker compose --profile dev up
```
Runs a Vite development server with live updates in port 8080

### Run Linting
```
docker compose --profile test run --rm frontend-test npm run lint
```
Runs biomes linter

### Run Tests
```
docker compose --profile test run --rm frontend-test npm run test
```
Runs the tests from the src/tests folder

### Run Build
```
docker compose --profile test run --rm frontend-test npm run build
```
Runs the build

### Run Production
```
docker compose --profile prod up
```
Runs a build of the application and provides the generated files over NGINX in port 8080