КАК ПОДНЯТЬ

Просто написать "make up" в корне. Если захочешь поменять порты - все в .env. Сейчас апи поднимается на порте 8080. В принципе, после make up можешь пойти попить чай или что такое - так как будет подниматься сервер и две базы данных (это не быстро, но только в первый раз, потом докер сохранит образы и будет быстро). Возможно, там постгрес с первого раза не успеет подняться, так что если увидешь, что не работает - сделай make down && make up.


КАК ПОЛЬЗОВАТЬСЯ

В данный реализовано:

- POST /api/v1/signup
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
}

- POST /api/v1/signin - в ответ будет вся инфа юзера!
Валидное тело:
{
	"email": "hello@gmail.com",
	"password": "123"
}
- POST /api/v1/user - обновление данных пользователя. Тут просто все данные пользователя обновятся на то, что отправлено. То есть нужно все эти поля обязательно присылать, если прислать какие - то пустые - они станут пустыми. Я думаю, так удобнее, так как на фронте все равно будут уже после signin все актуальные данные, и сюда ты просто их же присылаешь, изменяя те, которые поменял юзер. Я имею ввиду, просто реактивно выводишь на страничку, реактивно они обновляются у тебя, и если юзер что то сохраняет - отсылаешь такой запрос на обновление.ы 
Тело запроса:
{
	"id": "user id",
	"phone": "89671102000",
	"username": "Liza",
	"age": 17,
	"gender": "female",
	"country": "Russia",
	"city": "Moscow",
	"max_dist": 100,
	"look_for": "male",
	"min_age": 24,
	"max_age": 47
}

- DELETE /api/v1/signout - тут ничего, кроме самомого запроса отправлять не нужно, сервер просто сам обнулит сессию

- GET /api/v1/strangers (если придумаешь название получше - супер). Это основной запрос, выполняющийся после того, как человек залогинился. Он возвращает ему пачку пользователей сайта, подходяший по его критериям (возраст, город, пол). Типо лента в сайте знакомств.

- GET /api/v1/data/{id юзера} - получение данных юзера по его id (используется для загрузки данных пользоватетелей, поветивщих, или лайкнувших страницу юзера (там как раз будут их id)) - ответом будет вся информация пользователя, кроме его looked_by, liked_by, matched (так как эта инфа только ему принадлежит)

- POST /api/media/upload - загрука фоток. Тело:
{
    id: <id юзера>,
    isAvatar: Boolean (true можешь писать только, когда устанавливаешь аватар, при любом другом значении этого поля, даже при его отсутсвии, фотка просто сохраниться в галерею пользователя)
    user_image: <файл картинки (для того чтобы это был файл, надо просто указать тип инпута file)>
}

- GET /api/media/img/<id картинки> - получение картинки


ГРУППА:

- POST /api/v1/look
- POST /api/v1/like
- POST /api/v1/match
- в эти при запроса телом надо отпралять только одно поле "id":<id того пользователя, с которым совершается действие (соответственно, просмотр его страницы / лайк его страницы / матч с ним(про это расскажу подробнее))>
То есть, например, юзера зашел на странице другого юзера. Ты тогда с его браузера отправляет пост /api/v1/look с телом, содержащим id того пользователя, на страницу которого зашли. После этого юзера, который заходил на страницу, появится у него в поле looked_by.
То же самое с liked_by и matched. Но про это расскажу еще подробнее.

Пример тела, содежащего id:
{
	"id": "kjhaskfhs87yf9ay94hkhf2k298e"
}




Соответсвенно, все get запросы тебе никак не нужно валидировать - я там через бек вешаю уникальную куку на человека, и по ней потом понимаю от кого запрос.

Ясно, что перед тем, как тестировать strangers, нужно насоздавать чуть чуть юзеров (причем таких, чтобы друг другу подходили) - так как у тебя база будет пустая.

ЗАКЛЮЧИТЕЛЬНЫЕ ЗАМЕЧАНИЯ

Если какой то воспрос, или что то не работает - прям сразу звони или пиши, на свежую быстрее разберемся.

Проверяй чтобы соответсвовали все типы данных и названия полей при оправке на сервер (я там поменял born_date на age).

Там есть еще папка nginx - она пока не нужна и нигде не используется.

Запросы отправлять на http://localhost:8080

