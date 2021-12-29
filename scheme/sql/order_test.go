package sql

import "testing"

func Test_ConvertToOrder_Success(t *testing.T) {
	data := getTestOrderDataSuccess()

	sl := []map[string]string{
		data,
	}
	scheme := getCorrectlyOrderScheme()

	if res, err := ConvertToOrders(scheme, sl); err != nil || res == nil {
		t.Fatal(err)
	}

}

func getTestOrderDataSuccess() map[string]string {
	return map[string]string{
		"ref":                "0x80B9A4BF015829F711E9992701CAD5FF",
		"originId":           "01cad5ff-9927-11e9-80b9-a4bf015829f7",
		"docDate":            "28.06.4019",
		"name":               "Order 038856",
		"client":             "dfb9ed78-51a7-11e7-80fe-0cc47a5468a5",
		"docSum":             "2780.00",
		"internetOrderStage": "8",
		"dpp":                "0.0",
		"dppOD":              "0.0",
		"dppOS":              "0.0",
		"pickupPoint":        "Test pickup",
		"internetOrder":      "0x01",
		"sentSms":            "0x01",
		"orderType":          "1",
		"deliverySum":        "0.00",
		"deliveryWay":        "1",
		"wantedDateShipment": "28.06.4019",
		"extraInfo":          "test info",
		"comment":            "test comment",
		"agreement":          "test agreement",
		"stock":              "stock",
		"deliveryAddress":    " address",
		"deliveryArea":       "deliveryArear",
		"deliveryTimeFrom":   "01.01.2001T12:00:00",
		"deliveryTimeTo":     "01.01.2001T12:00:00",
		"doctor":             "Test doctor",
		"userId":             "475",
		"shipmentDate":       "28.06.4019",
		"department":         "department",
		"prepaid":            "74",
		"prepayment":         "2780.00",
		"credit":             "0.00",
	}
}

func getCorrectlyOrderScheme() []Field {
	return []Field{
		Field{Key: "originId", Required: true},
		Field{Key: "docDate", Required: true},
		Field{Key: "name", Required: true},
		Field{Key: "client", Required: true},
		Field{Key: "docSum", Required: true},
		Field{Key: "internetOrderStage", Required: true},
		Field{Key: "dpp", Required: true},
		Field{Key: "dppOD", Required: true},
		Field{Key: "dppOS", Required: true},
		Field{Key: "pickupPoint", Required: true},
		Field{Key: "internetOrder", Required: true},
		Field{Key: "sentSms", Required: true},
		Field{Key: "orderType", Required: true},
		Field{Key: "deliverySum", Required: true},
		Field{Key: "deliveryWay", Required: true},
		Field{Key: "wantedDateShipment", Required: true},
		Field{Key: "extraInfo", Required: true},
		Field{Key: "comment", Required: true},
		Field{Key: "agreement", Required: true},
		Field{Key: "stock", Required: true},
		Field{Key: "deliveryAddress", Required: true},
		Field{Key: "deliveryArea", Required: true},
		Field{Key: "deliveryTimeFrom", Required: true},
		Field{Key: "deliveryTimeTo", Required: true},
		Field{Key: "doctor", Required: true},
		Field{Key: "userId", Required: true},
		Field{Key: "shipmentDate", Required: true},
		Field{Key: "department", Required: true},
		Field{Key: "prepaid", Required: true},
		Field{Key: "prepayment", Required: true},
		Field{Key: "credit", Required: true},
	}
}
