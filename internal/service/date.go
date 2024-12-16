package service

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/paudarco/todo/internal/errors"
	"github.com/paudarco/todo/internal/models"
)

// NextDate возвращает следующую дату события, согласно заданному правилу.
func (s *TodoItemService) NextDate(now time.Time, date string, repeat string) (string, error) {
	// Парсим правило через функцию parseRule
	// repeatHeader - заголовок правила повторения (y, d, w, m).
	// repeatRule - словарь с ["days"] количеством дней для периода повторения, днями недели
	// или конкретными числами месяца (в зависимости от заголовка правила)
	// и ["months"] с номерами месяцев для определения месяцев повторения
	repeatHeader, repeatRule, err := parseRule(repeat)
	if err != nil {
		return "", err
	}

	// Парсим дату события.
	nextDate, err := time.Parse(models.Format, date)
	if err != nil {
		return "", errors.ErrInvalidDate
	}

	switch repeatHeader {
	case "y":
		for {
			// Увеличиваем дату ровно на год.
			nextDate = nextDate.AddDate(1, 0, 0)

			// Если дата раньше сегодняшнего дня - запускаем новую итерацию
			// увеличения года по бесконечному циклу.
			if nextDate.After(now) {
				return nextDate.Format(models.Format), nil
			}
		}
	case "d":
		for {
			// Увеличиваем дату на указанное колличество дней
			nextDate = nextDate.AddDate(0, 0, repeatRule["days"][0])

			// Если дата раньше сегодняшнего дня - запускаем новую итерацию
			// увеличения дней.
			if nextDate.After(now) {
				return nextDate.Format(models.Format), nil
			}
		}
	case "w":
		for {
			// Каждую итерацию цикла увеличиваем день на 1.
			nextDate = nextDate.AddDate(0, 0, 1)

			// Если день недели не соответствует любому дню из правила, продолжаем цикл.
			for _, day := range repeatRule["days"] {
				if day == int(nextDate.Weekday()) && nextDate.After(now) {
					return nextDate.Format(models.Format), nil
				}
			}
		}
	case "m":
		// В этом правиле есть два алгоритма: с указанными месяцами в правиле и без.
		if _, ok := repeatRule["months"]; ok {
			// Сразу прибавляем один день к сегодняшней дате,
			// чтобы она не участвовала в проверке.
			nextDate = nextDate.AddDate(0, 0, 1)

			for {
				for _, month := range repeatRule["months"] {
					// Если нынешний месяц соответствует любому и перечисленных в правиле,
					// начинаем искать нужный день, каждую итерацию увеличивая день на 1.
					for ; month == int(nextDate.Month()); nextDate = nextDate.AddDate(0, 0, 1) {
						for _, day := range repeatRule["days"] {
							// Если в правиле дня указано -1 или -2 - искомое число будет
							// последним или предпоследним относительно нынешнего месяца.
							// Для этого мы переходим к первому числу следующего месяца и вычитаем 1 или 2 дня.
							if day == -1 || day == -2 {
								day = time.Date(nextDate.Year(), nextDate.Month()+1, 1+day, 0, 0, 0, 0, time.UTC).Day()
							}

							if day == nextDate.Day() && nextDate.After(now) {
								return nextDate.Format(models.Format), nil
							}
						}
					}
				}

				// Если в списке месяцев в правиле нет текущего месяца,
				// переходим к первому числу следующего месяца.
				nextDate = time.Date(nextDate.Year(), nextDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)
			}
		}

		// Если правило не содержит месяцев - ищем день в текущем и последующих месяцах.
		for {
			// Каждую итерацию цикла увеличиваем день на 1,
			// пока не найдем искомый.
			nextDate = nextDate.AddDate(0, 0, 1)

			for _, day := range repeatRule["days"] {
				if day == -1 || day == -2 {
					day = time.Date(nextDate.Year(), nextDate.Month()+1, 1+day, 0, 0, 0, 0, time.UTC).Day()
				}

				if day == nextDate.Day() && nextDate.After(now) {
					return nextDate.Format(models.Format), nil
				}
			}
		}

	default:
		return "", errors.ErrInvalidRepeatRule
	}
}

func parseRule(repeat string) (string, map[string][]int, error) {
	// Если правило пустое или начинается с пробела - возвращаем ошибку
	if len(repeat) == 0 || strings.HasPrefix(repeat, " ") {
		return "", make(map[string][]int), errors.ErrInvalidRepeatRule
	}

	// Делим  правило на составные части и создаем результирующие переменные:
	// repeatRule - правила повторений по дням и месяцам
	// repeatHeader - заголовок правила повторений (y, d, w, m)
	repeatArr := strings.Split(repeat, " ")
	repeatRule := make(map[string][]int)
	var repeatHeader string

	switch repeatArr[0] {
	case "y":
		// В правиле не может быть больше одного символа.
		if len(repeatArr) != 1 {
			return "", repeatRule, errors.ErrInvalidRepeatRule
		}

		repeatHeader = repeatArr[0]

		return repeatHeader, repeatRule, nil

	case "d":
		// В правиле может быть только 2 части.
		if len(repeatArr) != 2 {
			return "", repeatRule, errors.ErrInvalidRepeatRule
		}

		// Проверяем, приводится ли вторая часть к целочисленному типу
		// и если да, соответствует ли она условиям.
		days, err := strconv.Atoi(repeatArr[1])
		if err != nil || days < 1 || days > 400 {
			return "", repeatRule, errors.ErrInvalidDay
		}

		repeatHeader = repeatArr[0]
		repeatRule["days"] = []int{days}

		return repeatHeader, repeatRule, nil

	case "w":
		// В правиле может быть только 2 части.
		if len(repeatArr) != 2 {
			return "", repeatRule, errors.ErrInvalidRepeatRule
		}

		// Мапа для проверки дней на повторение в правиле.
		days := make(map[int]bool)

		// Числа разделены запятой, поэтому указываем ее как разделитель
		// и проходимся по каждому элементу получившеголся слайса.
		for _, day := range strings.Split(repeatArr[1], ",") {
			dayNum, err := strconv.Atoi(day)
			// Дней в неделе - максимум 7
			if err != nil || dayNum < 1 || dayNum > 7 {
				return "", repeatRule, errors.ErrInvalidDay
			}

			// Если день повторился - возвращаем ошибку о неправильном правиле
			if days[dayNum] {
				return "", repeatRule, errors.ErrDublicateDay
			}

			days[dayNum] = true

			// В пакете time 7 день (воскресенье) имеет индекс 0,
			// т. к. неделя начинается именно с него.
			// Для соответствия документации пакета меняем порядковый номер дня.
			if dayNum == 7 {
				dayNum = 0
			}

			repeatRule["days"] = append(repeatRule["days"], dayNum)
		}

		repeatHeader = repeatArr[0]
		// Сортируем получившийся слайс номеров дней недели.
		slices.Sort(repeatRule["days"])
		return repeatHeader, repeatRule, nil

	case "m":
		// В правиле может быть либо 2, либо 3 части (дополнительная часть - месяц).
		if len(repeatArr) < 2 || len(repeatArr) > 3 {
			return "", repeatRule, errors.ErrInvalidRepeatRule
		}

		// Мапа для проверки дней на повторение в правиле.
		days := make(map[int]bool)

		for _, day := range strings.Split(repeatArr[1], ",") {
			dayNum, err := strconv.Atoi(day)
			// День имеет номер от 1 до 31
			// Также -1 и -2 указывают на последний и предпоследний день месяца.
			if err != nil || dayNum < -2 || dayNum > 31 || dayNum == 0 {
				return "", repeatRule, errors.ErrInvalidDay
			}

			if days[dayNum] {
				return "", repeatRule, errors.ErrDublicateDay
			}

			days[dayNum] = true

			repeatRule["days"] = append(repeatRule["days"], dayNum)
		}

		repeatHeader = repeatArr[0]
		slices.Sort(repeatRule["days"])

		// В отсортированном по порядку слайсе days -1 и -2 будут находиться
		// в самом начале, что логически неправильно, т.к. это последний и предпоследний дни.
		// Если они присутствуют в слайсе - по очереди помещаем их в самый конец.
		if repeatRule["days"][0] == -2 {
			repeatRule["days"] = append(repeatRule["days"][1:], repeatRule["days"][:1]...)
		}
		if repeatRule["days"][0] == -1 {
			repeatRule["days"] = append(repeatRule["days"][1:], repeatRule["days"][:1]...)
		}

		// Если правило повторения для месяца отсутствует - заканчиваем обработку правила.
		if len(repeatArr) == 2 {
			return repeatHeader, repeatRule, nil
		}

		// Мапа для проверки месяцев на повторение в правиле
		months := make(map[int]bool)

		for _, month := range strings.Split(repeatArr[2], ",") {
			monthNum, err := strconv.Atoi(month)

			if err != nil || monthNum < 1 || monthNum > 12 {
				return "", repeatRule, errors.ErrInvalidMonth
			}

			if months[monthNum] {
				return "", repeatRule, errors.ErrDublicateMonth
			}

			months[monthNum] = true
			repeatRule["months"] = append(repeatRule["months"], monthNum)
		}

		// Сортируем месяца по порядку.
		slices.Sort(repeatRule["months"])

		return repeatHeader, repeatRule, nil

	default:
		return "", repeatRule, errors.ErrInvalidRepeatRule
	}
}
