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
        x, h, a = map(float, input_values)
        break
    except ValueError:
        print("Error: введите числа")

R = calculate_r(n, a, x, h)
print(" ".join("{0:.5f}".format(num) for num in R))
