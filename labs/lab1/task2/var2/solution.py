max_value = max(R) if R else None  # Находим максимум, если список не пустой
if max_value is not None and max_value in R:
    max_index = R.index(max_value)
else:
    print("Notice: максимальный элемент не найден")

result = []
for i, x in enumerate(R):
    if i >= max_index or x >= 0:
        result.append(x)


print(" ".join(map(str, result)))
