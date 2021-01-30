Бэкенд GO - Тестовое задание
Необходимо реализовать сервис на языке Golang, который занимается парсингом, хранением и выдачей курсов валют.
Раз в 3 минуты этот сервис должен обновлять курсы валют в базе данных. Данные по курсам валют можно взять отсюда: http://www.cbr.ru/scripts/XML_daily.asp или https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml
Реализовать 2 REST API метода:
GET /currencies — должен возвращать список курсов валют с возможность пагинации
GET /currency/ — должен возвращать курс валюты для переданного id или кода валюты
API должно быть закрыто bearer авторизацией.
Дополнительно:
1. Хранить историю курсов за последний 30 минут
2. Написать конфигфайл, для запуска в докере

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