first_neg = next((i for i, x in enumerate(R) if x < 0), None)

if first_neg is None:
    print("Notice: все элементы положительные")
else:
    result = []
    for i, x in enumerate(R):
        if i <= first_neg or x <= 0:
            result.append(x)

    print(" ".join(map(str, result)))
