# news_api

## [Тестовое задание](https://gist.github.com/bethrezen/d6f17fbb039a4366fe6baafdf189ff9a)

Сделаны все основные задания + дополнительные

### Используемый стек

* **Golang 1.22**
* **Fiber** основной веб фреймворк
* **PostgreSQL** как основная БД
* **Reform** ORM для работы с БД
* **golang-migrate/migrate** для миграций основной бд
* **logrus** для логирования
* **Docker и Docker Compose** для быстрого развертывания


### Вопросы по тестовому заданию

В процессе разработки я столкнулся с некоторыми вопросами касательно некоторых моментов:

**Ручка /edit:id**  
В задании не указано, можно ли задать пустое значение входящего массива Categories (т.е. новость должна быть без каких-либо категорий)?
Решил, что **_нельзя задавать новость без категорий_**, поэтому если на вход идет пустой массив, то поле не обновляется

**Типы таблицы NewsCategories**  
В примере БД указано, что все 2 поля таблицы NewsCategories это первичные ключи (в postgresql в таблице может быть максимум 1 ключ). 
В задании не сказано каким образом ведется создание новостей (получается, что отсутствует какая-либо связь между таблицами News и NewsCategories, потому что во второй индексы создаются независимо)
Принято решение **_создать отдельный первичный ключ, исходные поля сделать обычными_**