package converter

import (
	"fmt"
	"strconv"
	"time"
)

func GetBitrixSegment(v string) (string, error) {
	segments := []int{
		0, 592, 593, 594, 595, 596, 597, 598, 599, 600, 601, 602, 603, 604, 605,
		606, 607, 608, 609, 610, 611, 612, 613, 614, 615, 616, 617, 618,
	}

	if i, err := strconv.Atoi(v); err == nil {
		return strconv.Itoa(segments[i]), nil
	} else {
		return "", err
	}
}

func GetCategoryInOrder(orderType string, internetOrder bool) int {
	var (
		offlineSales int = 7
		onlineSales  int = 4
		remake       int = 8
	)
	if internetOrder {
		return onlineSales
	} else if orderType == "1" {
		return offlineSales
	} else {
		return remake
	}
}

func GetInternetOrderStage(stage string) (string, error) {
	stages := map[string]string{
		"1":  "C4:NEW",
		"2":  "C4:5",
		"3":  "C4:PREPARATION",
		"4":  "C4:6",
		"5":  "C4:PREPAYMENT_INVOICE",
		"6":  "C4:EXECUTING",
		"7":  "C4:7",
		"8":  "C4:WON",
		"9":  "C4:LOSE",
		"10": "C4:8",
	}
	if s, ok := stages[stage]; !ok {
		return "", fmt.Errorf("can't get stage. Check map which is stages storage")
	} else {
		return s, nil
	}
}

func GetOrderStageByKind(kind string, closed bool) string {
	var pipline int
	var stage string
	if kind == "1" {
		pipline = 7
	} else {
		pipline = 8
	}

	if closed {
		stage = "WON"
	} else {
		stage = "NEW"
	}
	return fmt.Sprintf("C%d:%s", pipline, stage)

}

func SubtractionYearsOffset(dstr string, offset int, format string) string {
	if date, err := time.Parse(format, dstr); err == nil {
		return date.AddDate(-offset, 0, 0).Format(format)
	} else {
		return dstr
	}
}

func GetOrderType(value string) (res string) {
	switch value {
	case "1":
		res = "Первичный"
	case "2":
		res = "Гарантийный"
	case "3":
		res = "СервисноеОбслуживание"
	case "4":
		res = "Переделка"
	}
	return
}

func GetDeliveryType(value string) (res string) {
	switch value {
	case "1":
		res = "Самовывоз"
	case "2":
		res = "До клиента"
	case "3":
		res = "Силами перевозчика"
	case "4":
		res = "Силами перевозчика по адресу"
	}
	return
}

func GetNameBrandFieldNameForSegment(key string) (field string, ok bool) {
	switch key {
	case "1":
		return "UF_CRM_1640555139512", true
	case "2":
		return "UF_CRM_1640554426710", true
	case "3":
		return "UF_CRM_1640554667672", true
	case "4":
		return "UF_CRM_1640555386148", true
	case "5":
		return "UF_CRM_1640555668909", true
	case "6":
		return "UF_CRM_1640555702499", true
	case "7":
		return "UF_CRM_1640556314186", true
	case "8":
		return "UF_CRM_1640556671958", true
	case "9":
		return "UF_CRM_1640556915068", true
	case "10":
		return "UF_CRM_1640557078219", true
	case "13":
		return "UF_CRM_1640557336955", true
	case "14":
		return "UF_CRM_1640557493265", true
	case "18":
		return "UF_CRM_1640557859657", true
	case "19":
		return "UF_CRM_1640558121295", true
	default:
		return "", false
	}
}
