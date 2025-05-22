if not R:
    print("empty")
    exit()

max_index = R.index(max(R))
selected = R[max_index+1:]
avg = sum(selected)/len(selected) if selected else 0
print(avg)