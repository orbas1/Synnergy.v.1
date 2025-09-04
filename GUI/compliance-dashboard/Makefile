install:
npm ci

lint:
npx eslint . --ext .ts

build:
npm run build

test:
npm test

start:
npm start

docker-build:
docker build -t compliance-dashboard .

docker-run:
docker run -p 3000:3000 compliance-dashboard
