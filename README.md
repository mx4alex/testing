# testing
100% покрытие тестами поискового сервиса

## Описание сервиса
1. SearchClient - структура с методом FindUsers, который отправляет запрос во внешнюю систему и возвращает результат, немного преобразуя его.
2. SearchServer - своего рода внешняя система. Непосредственно занимается поиском данных в файле [dataset.xml](./dataset.xml). В продакшене бы запускалась в виде отдельного веб-сервиса.

- [coverage_test.go](./coverage_test.go) - файл с тестами
- [cover.html](./cover.html) - html-отчет с покрытием
- [client.go](./client.go) - клиент
- [server.go](./server.go) - сервер
