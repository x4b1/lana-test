package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	c := cli{baseUrl: os.Getenv("BASE_URL")}

	c.Run()
}

type cli struct {
	basket  string
	baseUrl string
}

func (c cli) Run() {
	prompt := promptui.Select{
		Label: "Select Action",
		Items: []string{
			"Create basket",
			"Add item",
			"Get total",
			"Delete basket",
		},
	}

	_, action, err := prompt.Run()
	if err != nil {
		fmt.Printf("Unexpected error %v\n", err)
		return
	}

	switch action {
	case "Create basket":
		c.CreateBasket()
	case "Add item":
		c.AddItem()
	case "Get total":
		c.GetTotal()
	case "Delete basket":
		c.DeleteBasket()
	}

	c.Run()
}

func (c *cli) CreateBasket() {
	res, err := http.Post(fmt.Sprintf("%s/baskets", c.baseUrl), "application/json", nil)
	if err != nil {
		fmt.Printf("Unexpected error %v\n", err)
		return
	}

	if res.StatusCode != http.StatusCreated {
		c.PrintError(res)
		return
	}

	var basket struct {
		ID string `json:"id"`
	}
	decodeResponse(res, &basket)

	c.basket = basket.ID

	fmt.Printf("%s basket created!!!\n", c.basket)
}

func (c *cli) AddItem() {
	prompt := promptui.Select{
		Label: "Select a product",
		Items: []string{
			"PEN",
			"TSHIRT",
			"MUG",
		},
	}

	_, product, err := prompt.Run()
	if err != nil {
		fmt.Printf("Unexpected error %v\n", err)
		return
	}

	res, err := http.Post(
		fmt.Sprintf("%s/baskets/%s/items", c.baseUrl, c.basket),
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"product":"%s"}`, product)),
	)
	if err != nil {
		fmt.Printf("Unexpected error %v\n", err)
		return
	}

	if res.StatusCode != http.StatusCreated {
		c.PrintError(res)
		return
	}

	fmt.Printf("%s added to %s basket!!!\n", product, c.basket)
}

func (c *cli) GetTotal() {
	res, err := http.Get(fmt.Sprintf("%s/baskets/%s/total", c.baseUrl, c.basket))
	if err != nil {
		fmt.Printf("Unexpected error %v\n", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		c.PrintError(res)
		return
	}

	var total struct {
		Total string `json:"total"`
	}
	decodeResponse(res, &total)

	fmt.Printf("total amount of the basket %s\n", total.Total)
}

func (c *cli) DeleteBasket() {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/baskets/%s", c.baseUrl, c.basket), nil)
	if err != nil {
		fmt.Printf("Unexpected error %v\n", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Unexpected error %v\n", err)
		return
	}

	if res.StatusCode != http.StatusNoContent {
		c.PrintError(res)
		return
	}

	var total struct {
		Total string `json:"total"`
	}
	decodeResponse(res, &total)

	fmt.Printf("Basket deleted %s\n", total.Total)
}

func (c cli) PrintError(res *http.Response) {
	var errorResponse struct {
		Message string `json:"message"`
	}
	decodeResponse(res, &errorResponse)

	fmt.Printf("Error from api code %s, message: %s\n", res.Status, errorResponse.Message)
}

func decodeResponse(r *http.Response, payload interface{}) error {
	return json.NewDecoder(r.Body).Decode(payload)
}
