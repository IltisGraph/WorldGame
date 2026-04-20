package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetCountries(c echo.Context) error {
	countries, err := gorm.G[PreviewCountry](DB).Where("visible = ?", true).Find(c.Request().Context())

	if err != nil {
		return err
	}

	return c.JSON(200, countries)
}
