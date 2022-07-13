package directus

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WebhookServices interface {
	FetchDataFromDirectusService(ctx context.Context) (string, error)
	CallDirectusService(ctx context.Context) (string, error)
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) FetchDataFromDirectusService(ctx context.Context) (string, error) {
	return "ok", nil
}

func (s *Service) CallDirectusService(ctx context.Context) (string, error) {

	url := "http://0.0.0.0:8055/items/consent_action_forms?fields=*,consent_id.*&filter[consent_id][status][_eq]=active&sort[]=-created_at"
	// url := "https://app-stg.fillgoods.co/stable/crm/items/consent_action_forms?fields=*,consent_id.*&filter[consent_id][is_active][_eq]=true&filter[is_active][_eq]=true&sort[]=created_at"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	req.Header.Add("Authorization", "Bearer test")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	fmt.Println("body", string(body))
	return "ok", nil
}
