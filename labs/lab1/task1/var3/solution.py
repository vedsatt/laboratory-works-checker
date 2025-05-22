import math

n = int(input())
x, h, a = map(float, input().split())

R = []
for i in range(1, n+1):
    R.append(6 * math.cos(a*x + i*h))

print(" ".join(map(str, R)))
