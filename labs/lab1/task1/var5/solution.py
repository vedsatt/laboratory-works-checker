import math
import sys
import io
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

def calculate_r(n, a, x, h):
    r = []
    for i in range(1, n + 1):
        value = math.cos(i * h) - math.cos(a * x + i * h)
        r.append(value)
    return r

while True:
    try:
        n = int(input(""))
        if n <= 0:
            print("Error: число должно быть > 0")
        elif n > 200:
            print("Error: число должно быть меньше 200")
        else:
            break
    except ValueError:
        print("Error: введите целое число")

while True:
    try:
        input_values = input().split()
        # Заменяем запятые на точки в каждом элементе списка
        input_values = [s.replace(',', '.') for s in input_values]
        x, h, a = map(float, input_values)
        break
    except ValueError:
        print("Error: введите числа в формате 8.5 или 8,5")

R = calculate_r(n, a, x, h)

# Форматируем числа и убираем лишние нули
formatted = []
for num in R:
    s = f"{num:.6f}"  # Форматируем с 6 знаками после запятой
    formatted.append(s)

# Выводим результат без лишних пробелов
print(" ".join(formatted))
