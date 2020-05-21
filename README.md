КАК ПОДНЯТЬ

Просто написать "make up" в корне. Если захочешь поменять порты - все в .env. Сейчас апи поднимается на порте 8080. В принципе, после make up можешь пойти попить чай или что такое - так как будет подниматься сервер и две базы данных (это не быстро, но только в первый раз, потом докер сохранит образы и будет быстро). Возможно, там постгрес с первого раза не успеет подняться, так что если увидешь, что не работает - сделай make down && make up.


КАК ПОЛЬЗОВАТЬСЯ

В данный момент реализовано 4 запроса:

1. POST /api/v1/signup
Валидное тело запроса:
{
	"email": "test3@gmail.com",
	"phone": "89671102000",
	"password": "123",
	"username": "Mary",
	"age": 20,
	"gender": "female",
	"country": "Russia",
	"city": "Moscow",
	"max_dist": 100,
	"look_for": "male",
	"min_age": 17,
	"max_age": 38,
	"images":["5eb23b466bfad500075db570","5eb23b596bfad500075db571","5eb23e4943e09d00075b7a22"]}
}

2. POST /api/v1/signin
Валидное тело:
{
	"email": "hello@gmail.com",
	"password": "123"
}
3. DELETE /api/v1/signout - тут ничего, кроме самомого запроса отправлять не нужно, сервер просто сам обнулит сессию
4. GET /api/v1/strangers (если придумаешь название получше - супер). Это основной запрос, выполняющийся после того, как человек залогинился. Он возвращает ему пачку пользователей сайта, подходяший по его критериям (возраст, город, пол). Типо лента в сайте знакомств.

5. POST /api/media/upload - загрука фоток. Тело:
{
    id: <id юзера>,
    isAvatar: Boolean
    user_image: <файл картинки (для того чтобы это был файл, надо просто указать тип инпута file)>
}

6. GET /api/media/img/<id картинки> - получение картинки

Соответсвенно, все get запросы тебе никак не нужно валидировать - я там через бек вешаю уникальную куку на человека, и по ней потом понимаю от кого запрос.

Ясно, что перед тем, как тестировать strangers, нужно насоздавать чуть чуть юзеров (причем таких, чтобы друг другу подходили) - так как у тебя база будет пустая.

ЗАКЛЮЧИТЕЛЬНЫЕ ЗАМЕЧАНИЯ

Если какой то воспрос, или что то не работает - прям сразу звони или пиши, на свежую быстрее разберемся.

Проверяй чтобы соответсвовали все типы данных и названия полей при оправке на сервер (я там поменял born_date на age).

Там есть еще папка nginx - она пока не нужна и нигде не используется.

Запросы отправлять на http://localhost:8080

