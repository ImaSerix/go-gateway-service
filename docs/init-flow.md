# Порядок выполнения инициализации сервисса

1. Получение файла конфигурации `config.yaml`
2. Чтение файла и преобразовывание его в объект `config.Root` (`LoadConfig` функция в [config.go](../internal/config/config.go))

3. Регистрация `Check`, `Middleware`, `Transformer`

4. Билд главного хэндлера
    1. Билд глобальных middleware
    2. Билд роутов
        1. Билд локальных middleware 
        2. Билд чеков