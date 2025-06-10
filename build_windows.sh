#!/bin/bash

# Устанавливаем переменные окружения для кросс-компиляции
export GOOS=windows
export GOARCH=amd64

# Создаем директорию для сборки, если её нет
mkdir -p build

# Компилируем программу
echo "Компиляция для Windows..."
go build -o build/intellivue.exe ./cmd/intellivue

# Проверяем результат
if [ $? -eq 0 ]; then
    echo "Сборка успешно завершена!"
    echo "Исполняемый файл создан: build/intellivue.exe"
else
    echo "Ошибка при сборке!"
    exit 1
fi 