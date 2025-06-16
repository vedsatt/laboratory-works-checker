import math

while True:
    try:
        n = int(input(""))
        if n <= 0:
            print("Error: число должно быть > 0")
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

R = []
for i in range(1, n+1):
    R.append(1.25 * math.sin(3*a*x - i*h))

if not R:
    print("Notice: массив пустой")
    exit()

print(" ".join(map(str, R)))
