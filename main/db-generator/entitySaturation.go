package main

import (
	"os"
	"fmt"
	schemeConf "github.com/vi-la-muerto/bx24sync/scheme/sql"
	"github.com/vi-la-muerto/bx24sync/sql"
)

func setClientReception(res []schemeConf.Reception, executeQuery sql.Execute) ([]schemeConf.Reception, error) {
	result := make([]schemeConf.Reception, 0, len(res))

	rd, err := os.OpenFile("client.json", os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	scheme, err := schemeConf.CreateScheme(rd)
	if err != nil {
		return nil, err
	}

	for _, reception := range res {

		id := fmt.Sprintf("0x%s", reception.ClientId)

		args := map[string]string{
			"${client}": id,
		}

		data, err := executeQuery(args, "client")

		if res, err := schemeConf.ConvertToClients(scheme, data); err == nil {
			if len(res) > 0 {
				reception.Client = res[0]
			}
		} else {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		result = append(result, reception)
	}
	return result, err
}

func setClientOrders(res []schemeConf.Order, executeQuery sql.Execute) ([]schemeConf.Order, error) {
	result := make([]schemeConf.Order, 0, len(res))

	rd, err := os.OpenFile("client.json", os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	scheme, err := schemeConf.CreateScheme(rd)
	if err != nil {
		return nil, err
	}

	for _, order := range res {

		id := fmt.Sprintf("0x%s", order.ClientId)

		args := map[string]string{
			"${client}": id,
		}

		data, err := executeQuery(args, "client")

		if res, err := schemeConf.ConvertToClients(scheme, data); err == nil {
			if len(res) > 0 {
				order.Client = res[0]
			}
		} else {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		result = append(result, order)
	}
	return result, err
}

func setClientShipment(res []schemeConf.Shipment, executeQuery sql.Execute) ([]schemeConf.Shipment, error) {
	result := make([]schemeConf.Shipment, 0, len(res))

	rd, err := os.OpenFile("client.json", os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	scheme, err := schemeConf.CreateScheme(rd)
	if err != nil {
		return nil, err
	}

	for _, shipment := range res {

		id := fmt.Sprintf("0x%s", shipment.ClientId)

		args := map[string]string{
			"${client}": id,
		}

		data, err := executeQuery(args, "client")

		if res, err := schemeConf.ConvertToClients(scheme, data); err == nil {
			if len(res) > 0 {
				shipment.Client = res[0]
			}
		} else {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		result = append(result, shipment)
	}
	return result, err
}

func setOrderSegments(res []schemeConf.Order, executeQuery sql.Execute) ([]schemeConf.Order, error) {
	result := make([]schemeConf.Order, 0, len(res))

	for _, order := range res {

		id := fmt.Sprintf("0x%s", order.Ref)

		args := map[string]string{
			"${order}": id,
		}

		if data, err := executeQuery(args, "order_segments"); err == nil {
			order.LoadSegments(data)
		} else {
			return nil, err
		}

		result = append(result, order)

	}
	return result, nil
}

func setShipmentSegments(res []schemeConf.Shipment, executeQuery sql.Execute) ([]schemeConf.Shipment, error) {
	result := make([]schemeConf.Shipment, 0, len(res))

	for _, shipment := range res {

		id := fmt.Sprintf("0x%s", shipment.Ref)

		args := map[string]string{
			"${shipment}": id,
		}

		if data, err := executeQuery(args, "shipment_segments"); err == nil {
			shipment.LoadSegments(data)
		} else {
			return nil, err
		}

		result = append(result, shipment)

	}
	return result, nil
}

func setReceptionPropertyes(res []schemeConf.Reception, executeQuery sql.Execute) ([]schemeConf.Reception, error) {

	result := make([]schemeConf.Reception, 0, len(res))

	rd, err := os.OpenFile("reception_propertyes.json", os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	scheme, err := schemeConf.CreateScheme(rd)
	if err != nil {
		return nil, err
	}

	for _, reception := range res {

		args := map[string]string{
			"${reception}": fmt.Sprintf("0x%s", reception.Id),
		}
		data, err := executeQuery(args, "reception_propertyes")

		if res, err := schemeConf.ConvertToAdditionalFields(scheme, data); err == nil {
			reception.AdditionnalFields = res
		} else {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		result = append(result, reception)
	}
	return result, err
}
