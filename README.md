<h1>Калькулятор</h1>

<h2>Краткое описание</h2>
<p>Данный сервис может использоваться для различных вычислений в режиме онлайн.</p>
<p>Поддерживаемые операции:</p>
<ol>
  <li>Сложение (+)</li>
  <li>Вычитание (-)</li>
  <li>Умножение (*)</li>
  <li>Деление (/)</li>
  <li>Операции приоритета (скобки)</li>
</ol>

<h2>Настройка проекта</h2>
<ol>
  <li>Клонируйте репозиторий.</li>
  <li>Установите библиотеку для чтения .env:
    <pre><code>go get github.com/joho/godotenv</code></pre>
  </li>
  <li>Создайте в корне проекта файл <code>.env</code> и добавьте в него настройки (пример заполнения):
    <ul>
      <li>TIME_ADDITION_MS=1</li>
      <li>TIME_SUBTRACTION_MS=1</li>
      <li>TIME_MULTIPLICATIONS_MS=1</li>
      <li>TIME_DIVISIONS_MS=1</li>
      <li>COMPUTING_POWER=2</li>
    </ul>
  </li>
  <li>Запустите проект:
    <pre><code>go run cmd/main.go</code></pre>
  </li>
</ol>

<h2>Как это работает</h2>
<p>При запуске приложения стартует сервер, который принимает запросы пользователей, а также агент, который опрашивает сервер на наличие новых задач. Если задачи есть, агент их обрабатывает и обновляет статус выражения.</p>
<h2>Схема приложения</h2>
<img scr="schema.jpeg">

<h2>Примеры запросов</h2>

<h2>UI</h2>
<p>В корне проекта находится HTML-страница <code>index.html</code>, с помощью которой можно удобно выполнять описанные ниже запросы.</p>

<h2>Сервер</h2>
<h3>Для запросов рекомендуется использовать программу Postman</h3>
<p>Сервер запускается и принимает запросы на базовый эндпоинт <code>/api/v1/calculate</code>. В теле запроса в поле <code>"expression"</code> укажите выражение, чтобы получить результат о созданной задаче.</p>
<p><strong>POST http://127.0.0.1:8080/api/v1/calculate</strong></p>
<p>Вход:</p>
<pre><code>{
    "expression": "выражение"
}</code></pre>
<p>На выходе:</p>
<pre><code>{
    "id": ID созданного выражения
}</code></pre>

<p>Чтобы увидеть созданные выражения, перейдите по эндпоинту <code>/api/v1/expressions</code>.</p>
<p><strong>GET http://127.0.0.1:8080/api/v1/expressions</strong></p>
<p>Вход:</p>
<pre><code>{}</code></pre>
<p>На выходе:</p>
<pre><code>{
    "expressions": [
        {
            "id": 0,
            "status": "queued",
            "result": "0"
        },
        {
            "id": 1,
            "status": "solved",
            "result": "0"
        }
    ]
}</code></pre>

<p>Чтобы получить конкретное выражение, используйте эндпоинт <code>/api/v1/expression/:id</code>.</p>
<p><strong>GET http://127.0.0.1:8080/api/v1/expression/1</strong></p>
<p>Вход:</p>
<pre><code>{}</code></pre>
<p>На выходе:</p>
<pre><code>{
    "expression": {
        "id": 1,
        "status": "solved",
        "result": "1.00"
    }
}</code></pre>

<h2>Агент</h2>
<p>Агент работает независимо от пользователя и опрашивает сервер в зависимости от числа воркеров.</p>
<p><strong>Получение задачи</strong></p>
<p><strong>GET http://127.0.0.1:8080/api/v1/internal/task</strong></p>
<p>Вход:</p>
<pre><code>{}</code></pre>
<p>На выходе:</p>
<pre><code>{
    "id": 1,
    "arg1": 1,
    "arg2": 1,
    "operation": 43,
    "operation_time": 1
}</code></pre>

<p><strong>Отправка результата</strong></p>
<p><strong>POST http://127.0.0.1:8080/api/v1/internal/task</strong></p>
<p>Вход:</p>
<pre><code>{
  "id": 1,
  "result": 2.00
}</code></pre>

<h2>Тестирование</h2>
<p>Для скриптов <code>calc.go</code> и <code>app.go</code> написаны тесты. Для их запуска перейдите в соответствующую директорию и выполните:</p>
<pre><code>go test -v</code></pre>

<h2>Разработчик</h2>
<ol>
  <li>Илья Савченко</li>
  <li>c4talyst@yandex.ru</li>
</ol>

<h2>Технологии</h2>
<ol>
  <li>Golang v1.20</li>
  <li>github.com/joho/godotenv</li>
</ol>