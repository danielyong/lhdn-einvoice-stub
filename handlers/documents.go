package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type lhdnDocumentTypeVersion struct {
	ID            int64     `json:"id" xml:"id"`
	Name          string    `json:"name" xml:"name"`
	Description   string    `json:"description" xml:"description"`
	ActiveFrom    time.Time `json:"activeFrom" xml:"activeFrom"`
	ActiveTo      time.Time `json:"activeTo" xml:"activeTo"`
	VersionNumber float32   `json:"versionNumber" xml:"versionNumber"`
	Status        string    `json:"status" xml:"status"`
}

type lhdnDocumentType struct {
	ID                   int64                     `json:"id" xml:"id"`
	InvoiceTypeCode      int64                     `json:"invoiceTypeCode" xml:"invoiceTypeCode"`
	Description          string                    `json:"description" xml:"description"`
	ActiveFrom           time.Time                 `json:"activeFrom" xml:"activeFrom"`
	ActiveTo             time.Time                 `json:"activeTo" xml:"activeTo"`
	DocumentTypeVersions []lhdnDocumentTypeVersion `json:"documentTypeVersions" xml:"documentTypeVersions"`
}

func GetDocumentTypes(ctx *gin.Context) {
	activeTo, _ := time.Parse("2006-01-02", "2030-01-01")
	invoiceDocumentTypeVersions := make([]lhdnDocumentTypeVersion, 0)
	invoiceDocumentTypeVersions = append(invoiceDocumentTypeVersions, lhdnDocumentTypeVersion{
		ID:            45,
		Description:   "invoice",
		ActiveFrom:    time.Time{},
		ActiveTo:      activeTo,
		VersionNumber: 1.0,
		Status:        "published",
	})

	somethingDocumentTypeVersions := make([]lhdnDocumentTypeVersion, 0)
	somethingDocumentTypeVersions = append(somethingDocumentTypeVersions, lhdnDocumentTypeVersion{
		ID:            69,
		Description:   "something",
		ActiveFrom:    time.Time{},
		ActiveTo:      activeTo,
		VersionNumber: 1.0,
		Status:        "published",
	})

	lhdnDocumentTypes := make([]lhdnDocumentType, 0)
	lhdnDocumentTypes = append(lhdnDocumentTypes, lhdnDocumentType{
		ID:                   45,
		InvoiceTypeCode:      04,
		Description:          "invoice",
		ActiveFrom:           time.Time{},
		ActiveTo:             activeTo,
		DocumentTypeVersions: invoiceDocumentTypeVersions,
	})
	lhdnDocumentTypes = append(lhdnDocumentTypes, lhdnDocumentType{
		ID:                   69,
		InvoiceTypeCode:      9,
		Description:          "something",
		ActiveFrom:           time.Time{},
		ActiveTo:             activeTo,
		DocumentTypeVersions: somethingDocumentTypeVersions,
	})

	ctx.JSON(200, &lhdnDocumentTypes)
}

func GetDocumentType(ctx *gin.Context) {
	documentTypeId := ctx.Param("document_type_id")

	activeTo, _ := time.Parse("2006-01-02", "2030-01-01")

	if id, _ := strconv.Atoi(documentTypeId); id == 45 {
		invoiceDocumentTypeVersions := make([]lhdnDocumentTypeVersion, 0)
		invoiceDocumentTypeVersions = append(invoiceDocumentTypeVersions, lhdnDocumentTypeVersion{
			ID:            45,
			Description:   "invoice",
			ActiveFrom:    time.Time{},
			ActiveTo:      activeTo,
			VersionNumber: 1.0,
			Status:        "published",
		})

		ctx.JSON(200, lhdnDocumentType{
			ID:                   45,
			InvoiceTypeCode:      04,
			Description:          "invoice",
			ActiveFrom:           time.Time{},
			ActiveTo:             activeTo,
			DocumentTypeVersions: invoiceDocumentTypeVersions,
		})
	} else {
		somethingDocumentTypeVersions := make([]lhdnDocumentTypeVersion, 0)
		somethingDocumentTypeVersions = append(somethingDocumentTypeVersions, lhdnDocumentTypeVersion{
			ID:            69,
			Description:   "something",
			ActiveFrom:    time.Time{},
			ActiveTo:      activeTo,
			VersionNumber: 1.0,
			Status:        "published",
		})

		ctx.JSON(200, &lhdnDocumentType{
			ID:                   69,
			InvoiceTypeCode:      9,
			Description:          "something",
			ActiveFrom:           time.Time{},
			ActiveTo:             activeTo,
			DocumentTypeVersions: somethingDocumentTypeVersions,
		})
	}
}
