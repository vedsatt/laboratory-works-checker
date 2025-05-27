R = [x for x in R if abs(x) >= 0.7]

if not R:
    print("Notice: массив пустой")
    exit()
    
print(" ".join(map(str, R)))
