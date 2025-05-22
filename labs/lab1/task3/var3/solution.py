if not R:
    print("empty")
    exit()

min_val = min(R)
last_min_index = len(R) - 1 - R[::-1].index(min_val)
selected = R[:last_min_index+1]
avg = sum(selected)/len(selected)
print(avg)