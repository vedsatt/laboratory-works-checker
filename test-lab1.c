#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <windows.h>

void cp() {
    SetConsoleCP(1251);
    SetConsoleOutputCP(1251);
}

// создаем массив и записываем туда значения элементов
float* CreateArray(int n, float x, float h, float a) {
    float* array = (float*)malloc(n * sizeof(float));

    for (int i = 1; i <= n; i++) {
        array[i-1] = cos(i * h) - cos(a * x + i * h);
    }

    return array;
}

int Delete(float** array, int* size) {
    int index_max;
    int index_pos = -1;
    float max_elem = -1;
    int len = 0;

    // находим последний положит. элемент и первый макс элемент по модулю
    for (int i = 0; i < *size; i++) {
        if ((*array)[i] > 0) {
            index_pos = i;
        }
        if (fabs((*array)[i]) > max_elem) {
            max_elem = fabs((*array)[i]);
            index_max = i;
        }
    }

    // если нет положительного числа
    if (index_pos == -1) {
        return -1;
    }
    // находим мин и макс индексы элементов
    int start = index_max < index_pos ? index_max : index_pos;
    int end = index_max > index_pos ? index_max : index_pos;

    // если индексы совпадают или находятся рядом
    if ((end - start) < 2){
        return -2;
    }

    // удаляем элементы в ранее найденных границах
    for (int i = start + 1; i < (*size - (end - start - 1)); i++) {
        (*array)[i] = (*array)[i + (end - start - 1)];
    }
    *size = *size - (end - start - 1);
    *array = (float*)realloc(*array, *size * sizeof(float));

    return 0;
}

float ArithmeticMean(float* array, int size) {
    float min_elem = 10;
    float max_elem = -10;
    int index_min, index_max;;
    float sum = 0.0;
    int cnt = 0;

    // находим границы по мин числу и по макс числу
    for (int i = 0; i < size; i++) {
        if (fabs(array[i]) < min_elem) {
            min_elem = fabs(array[i]);
            index_min = i;
        }
        if (fabs(array[i]) > max_elem) {
            max_elem = fabs(array[i]);
            index_max = i;
        }
    }

    // проверяем координаты мин и макс элементов
    int start = index_min < index_max ? index_min : index_max;
    int end = index_min > index_max ? index_min : index_max;

    // считаем кол-во элементов и их сумму
    for (int i = start + 1; i < end; i++) {
        sum += array[i];
        cnt++;
    }

    return (cnt > 0) ? (sum / cnt) : NAN;
}

int main() {
    cp();
    float x, h, a;
    int n, check;
    float arifm;
    int limit = 200;

    printf("#> Лабораторная работа №1\n");

    printf("\n#> Задание 1\n");
    printf("#> Введите длину массива R от 1 до %d:\n", limit);

        do {
        check = scanf("%d", &n);  // scanf возвращает 1, если введено число
        
        // Если ввод не число (check != 1)
        if (check != 1) {
            printf("Error: Введено не число! Пожалуйста, введите целое число.\n");
            
            // Очистка буфера ввода, чтобы избежать зацикливания
            while (getchar() != '\n');
            continue;  // Пропускаем остальные проверки и запрашиваем ввод снова
        }
        
        // Проверка диапазона
        if (n > limit) {
            printf("Error: Число превышает лимит (%d)\n", limit);
        } else if (n <= 0) {
            printf("Error: Число должно быть больше 0\n");
        }
    } while (n > limit || n <= 0);

    printf("#> Введите значения x, h, a:\n");
    scanf("%f %f %f", &a, &h, &x);
    float* array = CreateArray(n, x, h, a);
    printf("#> R массив из %d элементов\n", n);
    for (int i = 0; i < n; i++) {
        printf("%.5f ", array[i]);
    }
    printf("\n\n#> Задание 2\n");
    check = Delete(&array, &n);
    switch (check) {
        case -1:
            printf("Notice: Нет положительного элемента\n");
            break;
        case -2:
            printf("Notice: Последний положительный и первый максимальный элементы находятся рядом или совпадают\n");
            break;
        default:
            printf("#> Массив R состоит из %d элементов\n", n);
            for (int i = 0; i < n; i++) {
                printf("%.5f ", array[i]);
            }
            printf("\n");
            break;
    }

    printf("\n#> Задание 3\n");
    arifm = ArithmeticMean(array, n);
    if (isnan(arifm)) {
        printf("Notice: Нет среднего значения\n");
    } else {
        printf("#> sr:\n%.5f\n", arifm);
    }

    free(array);
    return 0;
}