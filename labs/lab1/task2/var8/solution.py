def remove_elements(r):
    if not r:
        return [], 1
    
    # Находим последний положительный элемент
    last_positive = -1
    for i in range(len(r) - 1, -1, -1):
        if r[i] > 0:
            last_positive = i
            break
    
    if last_positive == -1:
        print("Notice: в массиве нет положительных элементов")
        return [], 1
    
    # Находим первый максимальный по модулю элемент
    max_abs_index = 0
    max_abs_value = abs(r[0])
    for i in range(1, len(r)):
        if abs(r[i]) > max_abs_value:
            max_abs_value = abs(r[i])
            max_abs_index = i
    
    # Проверяем, находятся ли элементы рядом или совпадают
    if abs(last_positive - max_abs_index) <= 1:
        print("Notice: последний положительный и первый максимальный по модулю элементы находятся рядом или совпадают")
        return [], 1
    
    # Удаляем элементы между последним положительным и первым максимальным по модулю
    start = min(last_positive, max_abs_index) + 1
    end = max(last_positive, max_abs_index)
    
    if start >= end:
        print("Notice: нет элементов для удаления между указанными позициями")
        return [], 1
    
    r = r[:start] + r[end:]
    return r, 0

R_filtered, flag = remove_elements(R.copy())
if len(R) == 0:
    print("Notice: пустой массив")
    exit()
elif not flag:
    print(" ".join("{0:.5f}".format(num) for num in R_filtered))


