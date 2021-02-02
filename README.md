#  Как запустить

Собрать приложение
```
docker-compose run --rm currency-loader-build
```

Собрать docker образ
```
docker-compose build currency-loader 
```

Запустить docker образ
```
docker-compose run -p 5000:5000 --name currency-loader --rm currency-loader
```

Остановить контейнер
```
docker stop  currency-loader
```

Теперь можно выполянть запросы

```
curl -X GET \
  'http://localhost:5000/currencies?page=2' \
  -H 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjN9.PZLMJBT9OIVG2qgp9hQr685oVYFgRgWpcSPmNcw6y7M' 
```
```
  curl -X GET \
  http://localhost:5000/currency/usd \
  -H 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjN9.PZLMJBT9OIVG2qgp9hQr685oVYFgRgWpcSPmNcw6y7M'
  ```