max_index = R.index(max(R))
result = []
for i, x in enumerate(R):
    if i >= max_index or x >= 0:
        result.append(x)

print(" ".join(map(str, result)))
