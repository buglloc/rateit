package contactsys

import (
	"fmt"
	"strings"
	"time"

	"github.com/jaswdr/faker"

	"github.com/buglloc/rateit/internal/phonearea"
)

type Transaction struct {
	ID string `json:"id"`
}

type TransferDetails struct {
	Amount   string `json:"trnAmount"`
	Currency string `json:"trnCurrency"`
	Source   string `json:"tSource"`
	Ground   string `json:"tGround"`
	Relation string `json:"tRelation"`
}

func RandomTransferDetails(amount float64) *TransferDetails {
	return &TransferDetails{
		Amount:   fmt.Sprintf("%.2f", amount),
		Currency: "USD",
		Source:   "1",
		Ground:   "1",
		Relation: "Wife",
	}
}

type TransferParticipants struct {
	*TransferSender
	*TransferRecipient
}

type TransferSender struct {
	Name         string `json:"sName"`
	LastName     string `json:"sLastName"`
	BirthPlace   string `json:"sBirthPlace"`
	Country      string `json:"sCountryC"`
	IDType       string `json:"sIDtype"`
	IDNumber     string `json:"sIDnumber"`
	IDExpireDate string `json:"sIDexpireDate"`
	Birthday     string `json:"sBirthday"`
	Sex          string `json:"sSex"`
	Occupation   string `json:"sOccupation"`
}

func RandomTransferSender() *TransferSender {
	fake := faker.New()
	person := fake.Person()
	return &TransferSender{
		Name:       safeName(person.FirstNameMale()),
		LastName:   safeName(person.LastName()),
		BirthPlace: "USSR",
		Country:    "RU",
		IDType:     "ПАСПОРТ ГРАЖДАНИНА РФ",
		IDNumber:   fake.Numerify("2715######"),
		IDExpireDate: time.Now().AddDate(
			fake.IntBetween(1, 5),
			fake.IntBetween(1, 12),
			fake.IntBetween(1, 20),
		).Format("2006-01-02"),
		Birthday: time.Now().AddDate(
			-fake.IntBetween(30, 50),
			fake.IntBetween(1, 12),
			fake.IntBetween(1, 20),
		).Format("2006-01-02"),
		Sex:        person.GenderMale(),
		Occupation: "Specialist",
	}
}

type TransferRecipient struct {
	Name     string `json:"bName"`
	LastName string `json:"bLastName"`
	Country  string `json:"bCountryC"`
	City     string `json:"bCity"`
	Phone    string `json:"bPhone"`
	Address  string `json:"bAddress"`
	Sex      string `json:"bSex"`
	Account  string `json:"bAccount"`
}

func RandomTransferRecipient(country Country, bank Bank) *TransferRecipient {
	fake := faker.New()
	person := fake.Person()
	address := fake.Address()
	return &TransferRecipient{
		Name:     safeName(person.FirstNameFemale()),
		LastName: safeName(person.LastName()),
		Country:  country.Code,
		City:     bank.Address,
		Phone:    fake.Numerify(phonearea.PhoneTemplate(country.Code)),
		Address: safeName(fmt.Sprintf(
			"%s %s %s",
			fake.Numerify("###"), address.StreetName(), address.SecondaryAddress(),
		)),
		Sex:     safeName(person.GenderFemale()),
		Account: fake.Numerify("#########"),
	}
}

func safeName(in string) string {
	return strings.ReplaceAll(in, `"`, "'")
}
