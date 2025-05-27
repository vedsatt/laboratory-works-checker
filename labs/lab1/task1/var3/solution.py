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
        x, h, a = map(float, input_values)
        break
    except ValueError:
        print("Error: введите числа")

R = []
for i in range(1, n + 1):
    R.append(6 * math.cos(a * x + i * h))

if not R:
    print("Notice: массив пустой")
    exit()

print(" ".join(map(str, R)))