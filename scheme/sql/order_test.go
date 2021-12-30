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
		{Key: "originId", Required: true},
		{Key: "docDate", Required: true},
		{Key: "name", Required: true},
		{Key: "client", Required: true},
		{Key: "docSum", Required: true},
		{Key: "internetOrderStage", Required: true},
		{Key: "dpp", Required: true},
		{Key: "dppOD", Required: true},
		{Key: "dppOS", Required: true},
		{Key: "pickupPoint", Required: true},
		{Key: "internetOrder", Required: true},
		{Key: "sentSms", Required: true},
		{Key: "orderType", Required: true},
		{Key: "deliverySum", Required: true},
		{Key: "deliveryWay", Required: true},
		{Key: "wantedDateShipment", Required: true},
		{Key: "extraInfo", Required: true},
		{Key: "comment", Required: true},
		{Key: "agreement", Required: true},
		{Key: "stock", Required: true},
		{Key: "deliveryAddress", Required: true},
		{Key: "deliveryArea", Required: true},
		{Key: "deliveryTimeFrom", Required: true},
		{Key: "deliveryTimeTo", Required: true},
		{Key: "doctor", Required: true},
		{Key: "userId", Required: true},
		{Key: "shipmentDate", Required: true},
		{Key: "department", Required: true},
		{Key: "prepaid", Required: true},
		{Key: "prepayment", Required: true},
		{Key: "credit", Required: true},
	}
}
