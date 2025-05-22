if not R:
    print("empty")
    exit()

min_index = R.index(min(R))
first_neg = next((i for i, x in enumerate(R) if x < 0), None)

if first_neg is None:
    print()
else:
    start = min(min_index, first_neg)
    end = max(min_index, first_neg)
    selected = R[start+1:end]
    avg = sum(selected)/len(selected) if selected else 0
    print(avg)