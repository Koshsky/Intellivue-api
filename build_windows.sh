#!/bin/bash

# Устанавливаем переменные окружения для кросс-компиляции
export GOOS=windows
export GOARCH=amd64

# Создаем директорию для сборки, если её нет
mkdir -p build

# Компилируем основное приложение
echo "Компиляция intellivue для Windows..."
go build -o build/intellivue.exe ./cmd/intellivue

# Проверяем результат компиляции основного приложения
if [ $? -ne 0 ]; then
    echo "Ошибка при сборке intellivue!"
    exit 1
fi

# Компилируем ресивер
echo "Компиляция receiver для Windows..."
go build -o build/receiver.exe ./cmd/receiver

# Проверяем результат компиляции ресивера
if [ $? -ne 0 ]; then
    echo "Ошибка при сборке receiver!"
    exit 1
fi

echo "Сборка успешно завершена!"
echo "Исполняемые файлы созданы:"
echo "- build/intellivue.exe"
echo "- build/receiver.exe" 