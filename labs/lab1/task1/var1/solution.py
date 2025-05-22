import math

n = int(input())
x, h, a = map(float, input().split())

R = []
for i in range(1, n+1):
    R.append(2.5 * math.sin(a*x + i**2 * h))

if not R:
    print("empty")
    exit()

print(" ".join(map(str, R)))
