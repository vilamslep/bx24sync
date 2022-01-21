package bitrix24

import (
	"strings"
	"testing"
)

func Test_NewDealFromJson(t *testing.T) {
	t.Fail()
}

func Test_GetDealFromRawAsReception_WithoutError(t *testing.T) {
	raw := receptionRaw()

	rd := strings.NewReader(raw)

	if _, err := GetDealFromRawAsReception(rd); err != nil {
		t.Fatalf("Can't to convert reception to deal")
	}
}
func Test_GetDealFromRawAsOrder_WithoutError(t *testing.T) {
	raw := orderRaw()

	rd := strings.NewReader(raw)

	if _, err := GetDealFromRawAsOrder(rd); err != nil {
		t.Fatalf("Can't to convert order to deal")
	}
}
func Test_GetDealFromRawAsShipment_WithoutError(t *testing.T) {
	raw := shipmentRaw()

	rd := strings.NewReader(raw)

	if _, err := GetDealFromRawAsShipment(rd); err != nil {
		t.Fatalf("Can't to convert shipment to deal")
	}
}
func Test_DealJson(t *testing.T) {
	t.Fail()
}
func Test_DealFind(t *testing.T) {
	t.Fail()
}
func Test_DealAdd(t *testing.T) {
	t.Fail()
}
func Test_DealUpdate(t *testing.T) {
	t.Fail()
}
func Test_checkContact(t *testing.T) {
	t.Fail()
}
func Test_getDealReader(t *testing.T) {
	t.Fail()
}

func receptionRaw() string {
	return `[{"id":"test_id","originId":"test_originID","name":"001159254","date":"17.11.2021","department":"department","client":"client","clientData":{"originId":"","name":"originId","birthday":"03.07.4013","gender":"1","isClient":"\u0001","isSuppler":"\u0000","otherRelation":"\u0000","isRetireeOrDisabledPerson":"\u0000","connectionWay":"100","thereIsContract":"\u0001","sendAds":"\u0000","isInternetClient":"\u0000","isOfflineClient":"\u0000","isClinicClient":"\u0001","discountClinicService":"3","discountMedicalThings":"5","discountRayban":"0","phone":"81111111111;","email":""},"userId":"219","usersFields":[{"Key":"UF_CRM_1580476536","Value":"03.07.2013"},{"Key":"UF_CRM_1589894793184","Value":""},{"Key":"UF_CRM_1589971684","Value":""},{"Key":"UF_CRM_1589972819","Value":""},{"Key":"UF_CRM_1589974106","Value":"\u0000"},{"Key":"UF_CRM_1589974192","Value":"\u0000"},{"Key":"UF_CRM_1589974192","Value":"\u0001"},{"Key":"UF_CRM_1589974231","Value":"\u0000"},{"Key":"UF_CRM_1589974231","Value":"\u0001"},{"Key":"UF_CRM_1589974273","Value":"\u0000"},{"Key":"UF_CRM_1589974338","Value":"\u0000"},{"Key":"UF_CRM_1589974338","Value":"\u0001"},{"Key":"UF_CRM_1589974369","Value":"\u0000"},{"Key":"UF_CRM_1589974403","Value":"\u0000"},{"Key":"UF_CRM_1589974403","Value":"\u0001"},{"Key":"UF_CRM_1589974471","Value":"\u0000"},{"Key":"UF_CRM_1589974504","Value":"\u0000"},{"Key":"UF_CRM_1589974545","Value":"\u0000"},{"Key":"UF_CRM_1589974668","Value":"\u0000"},{"Key":"UF_CRM_1589974668","Value":"\u0001"},{"Key":"UF_CRM_1589974697","Value":"\u0000"},{"Key":"UF_CRM_1589974731","Value":"\u0000"},{"Key":"UF_CRM_1589974767","Value":"\u0000"},{"Key":"UF_CRM_1589975534","Value":"+0,75"},{"Key":"UF_CRM_1589975707","Value":"+0,75"},{"Key":"UF_CRM_1589975793","Value":"-0,75"},{"Key":"UF_CRM_1589975820","Value":"171"},{"Key":"UF_CRM_1589975905","Value":"27,5"},{"Key":"UF_CRM_1589975966","Value":"54"},{"Key":"UF_CRM_1589975996","Value":"26,5"},{"Key":"UF_CRM_1589976546","Value":"8.0000"},{"Key":"UF_CRM_1590047492","Value":"100.0000"},{"Key":"UF_CRM_1590053284","Value":"17.11.2021"},{"Key":"UF_CRM_1590738041","Value":"81111111111"},{"Key":"UF_CRM_1591092748","Value":"1,0"},{"Key":"UF_CRM_1591092956","Value":"1,0"},{"Key":"UF_CRM_1591352675","Value":"Для постоянного ношения"}]}]`
}

func orderRaw() string {
	return `[{"ref":"test_Ref","originId":"test_OriginId","docDate":"17.11.4021","name":"ОМ0Т-084797","client":"test_client","clientData":{"originId":"test_ID","name":"test_Name","birthday":"18.08.3982","gender":"1","isClient":"\u0001","isSuppler":"\u0000","otherRelation":"\u0000","isRetireeOrDisabledPerson":"\u0000","connectionWay":"100","thereIsContract":"\u0000","sendAds":"\u0000","isInternetClient":"\u0000","isOfflineClient":"\u0001","isClinicClient":"\u0000","discountClinicService":"7","discountMedicalThings":"10","discountRayban":"0","phone":"81111111111;81111111111;","email":"main@main.ru"},"docSum":"560.00","internetOrderStage":"1","dpp":"0.0","dppOD":"0.0","dppOS":"0.0","pickUpPoint":"test_pickUp","intenterOrder":"","sentSms":"\u0000","orderType":"1","deliverySum":"0.00","deliveryWay":"1","wantedDateShipment":"17.11.4021","extraInfo":"","comment":"","agreement":"test_agreement","stock":"test_stock","deliveryAddress":"","deliveryArea":"","deliveryTimeFrom":"01.01.2001T12:00:00","deliveryTimeTo":"01.01.2001T12:00:00","doctor":"","userId":"110","shipmentDate":"17.11.4021","prepaid":"0.00","prepayment":"560.00","credit":"0.00","department":"test_department","segment":[{"id":"15","brands":["","",""]}],"closed":"\u0001","finishDate":"17.11.4021"}]`
}

func shipmentRaw() string {
	return `[{"ref":"ref","originId":"originId","name":"ОМ0Т-157006","docDate":"17.11.4021","client":"client","clientData":{"originId":"originId","name":"name","birthday":"29.11.4001","gender":"1","isClient":"\u0001","isSuppler":"\u0000","otherRelation":"\u0000","isRetireeOrDisabledPerson":"\u0000","connectionWay":"100","thereIsContract":"\u0001","sendAds":"\u0000","isInternetClient":"\u0000","isOfflineClient":"\u0000","isClinicClient":"\u0001","discountClinicService":"9","discountMedicalThings":"12","discountRayban":"0","phone":"81111111111;","email":""},"docSum":"1940.00","department":"department","stock":"stock","agreement":"agreement","comment":"comment","doctor":"","userId":"328","segment":[{"id":"18","brands":["","Biotrue 120 ml"]},{"id":"7","brands":["Premio"]}]}]`
}
