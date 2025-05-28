def calculate_average(r):
    if not r:
        return None
    
    # Находим первый максимальный элемент
    max_index = 0
    max_value = r[0]
    for i in range(1, len(r)):
        if r[i] > max_value:
            max_value = r[i]
            max_index = i
    
    # Находим минимальный по модулю элемент
    min_abs_index = 0
    min_abs_value = abs(r[0])
    for i in range(1, len(r)):
        if abs(r[i]) < min_abs_value:
            min_abs_value = abs(r[i])
            min_abs_index = i
    
    # Определяем границы для подсчета среднего
    start = min(max_index, min_abs_index) + 1
    end = max(max_index, min_abs_index)
    
    if start >= end:
        return None
    
    selected_elements = r[start:end]
    average = sum(selected_elements) / len(selected_elements) if selected_elements else None
    return average

average = calculate_average(R)
if average is None:
    print("Notice: нет ср. арифм. числа")
else:
    print("{0:.5f}".format(average))
