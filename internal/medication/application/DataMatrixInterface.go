package application

import (
	dataMatrix "github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/dataMatrix"
)

type DataMatrixCache interface {
	Get(key string) (dataMatrix.MedicationInfoFromAPI, error) // по идее импортировать dataMatrix нельзя а нужно весь этот интерфейс в домен запихать и там сделать MedicationInfoFromAPI
	Set(key string, data dataMatrix.MedicationInfoFromAPI) (error)
}

type DataMatrixClient interface {
	GetInformationByDataMatrix(GTIN string, SerialNumber string, CryptoData91 string, CryptoData92 string) (dataMatrix.MedicationInfoFromAPI, error)
}
