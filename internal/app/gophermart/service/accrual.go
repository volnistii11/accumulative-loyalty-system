package service

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"net/http"
	"strings"
)

type Accrual struct {
}

func NewAccrual() *Accrual {
	return &Accrual{}
}

type NewOrderGetter interface {
	GetNewOrders() []string
}

func (a *Accrual) GetNewOrders(db NewOrderGetter) []string {
	return db.GetNewOrders()
}

func (a *Accrual) SendOrderNumbersToAccrualSystem(orderNumber string, endpoint string) (*model.AccrualSystemAnswer, error) {
	//endpoint example: http://localhost:8080/api/orders/
	fmt.Println("SendOrderNumbersToAccrualSystem: start")
	endpointWithOrderNumber := fmt.Sprintf("%s%s", endpoint, orderNumber)
	client := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, endpointWithOrderNumber, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	defer response.Body.Close()
	var accrualSystemAnswer model.AccrualSystemAnswer
	err = json.NewDecoder(response.Body).Decode(&accrualSystemAnswer)
	fmt.Println(accrualSystemAnswer)
	if err != nil {
		return nil, err
	}

	return &accrualSystemAnswer, nil
}

type accrualInfoUpdater interface {
	UpdateAccrualInfoForOrderNumber(newInfo *model.AccrualSystemAnswer) error
}

func (a *Accrual) UpdateAccrualInfoForOrderNumber(db accrualInfoUpdater, newInfo *model.AccrualSystemAnswer) error {
	err := db.UpdateAccrualInfoForOrderNumber(newInfo)
	if err != nil {
		return err
	}
	return nil
}
