package handlers

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type submissionFormat struct {
	Documents []document `json:"documents" xml:"documents"`
}

type document struct {
	Format       string `json:"format" xml:"format"`
	Document     string `json:"document" xml:"document"`
	DocumentHash string `json:"documentHash" xml:"documentHash"`
	CodeNumber   string `json:"codeNumber" xml:"codeNumber"`
}

type acceptedDocument struct {
	UUID              string `json:"uuid" xml:"uuid" gorm:"primaryKey"`
	InvoiceCodeNumber string `json:"invoiceCodeNumber" xml:"invoiceCodeNumber"`
}

type acceptedDocumentExtended struct {
	UUID                  string    `json:"uuid" xml:"uuid" gorm:"primaryKey"`
	SubmissionUID         string    `json:"submissionUid" xml:"submissionUid"`
	LongID                string    `json:"longId" xml:"longId"`
	InvoiceCodeNumber     string    `json:"internalId" xml:"internalId"`
	TypeName              string    `json:"typeName" xml:"typeName"`
	TypeVersionName       string    `json:"typeVersionName" xml:"typeVersionName"`
	IssuerTIN             string    `json:"issuerTin" xml:"issuerTin"`
	ReceiverID            string    `json:"receiverId" xml:"receiverId"`
	ReceiverName          string    `json:"receiverName" xml:"receiverName"`
	DateTimeIssued        time.Time `json:"dateTimeIssued" xml:"dateTimeIssued"`
	DateTimeReceived      time.Time `json:"dateTimeReceived" xml:"dateTimeReceived"`
	DateTimeValidated     time.Time `json:"dateTimeValidated" xml:"dateTimeValidated"`
	TotalSales            float32   `json:"totalSales" xml:"totalSales"`
	TotalDiscount         float32   `json:"totalDiscount" xml:"totalDiscount"`
	NetAmount             float32   `json:"netAmount" xml:"netAmount"`
	Total                 float32   `json:"total" xml:"total"`
	Status                string    `json:"status" xml:"status"`
	CancelDateTime        time.Time `json:"cancelDateTime" xml:"cancelDateTime"`
	RejectRequestDateTime time.Time `json:"rejectRequestDateTime" xml:"rejectRequestDateTime"`
	DocumentStatusReason  string    `json:"documentStatusReason" xml:"documentStatusReason"`
	CreatedByUserID       string    `json:"createdByUserId" xml:"createdByUserId"`
}

type rejectedDocument struct {
	InvoiceCodeNumber string    `json:"invoiceCodeNumber" xml:"invoiceCodeNumber"`
	Error             lhdnError `json:"error" xml:"error"`
}

type lhdnError struct {
	Code    string      `json:"code" xml:"code"`
	Message string      `json:"message" xml:"message"`
	Target  string      `json:"target" xml:"target"`
	Details []lhdnError `json:"details" xml:"details"`
}

type outputDocument struct {
	SubmissionUID     string             `json:"submissionUID" xml:"submissionUID" gorm:"column:submission_uid"`
	AcceptedDocuments []acceptedDocument `json:"acceptedDocuments" xml:"acceptedDocuments"`
	RejectedDocuments []rejectedDocument `json:"rejectedDocuments" xml:"rejectedDocuments"`
}

type submission struct {
	SubmissionUID    string                     `json:"submissionUid"`
	DocumentCount    int64                      `json:"documentCount"`
	DateTimeReceived time.Time                  `json:"dateTimeReceived"`
	OverallStatus    string                     `json:"overallStatus"`
	DocumentSummary  []acceptedDocumentExtended `json:"documentSummary"`
}

type arbitraryXML struct {
	XMLName  xml.Name
	Attrs    []xml.Attr     `xml:",any,attr"`
	Content  string         `xml:",innerxml"`
	Children []arbitraryXML `xml:",any"`
}

func ReallySelectStarSubmission(ctx *gin.Context) {

	var documents []acceptedDocumentExtended

	if results := db.Find(&documents); results.RowsAffected > 0 {
		ctx.JSON(200, &documents)
	} else {
		ctx.JSON(400, &lhdnError{
			Code:    "BadArgument",
			Message: "No submissions received previously",
		})
	}
}

func SelectStarSubmission(ctx *gin.Context) {
	submissionUID := ctx.Param("submission_id")

	var documents []acceptedDocumentExtended

	if results := db.Find(&documents, "submission_uid = ?", submissionUID); results.RowsAffected > 0 {
		ctx.JSON(200, &submission{
			SubmissionUID:    submissionUID,
			DocumentCount:    int64(len(documents)),
			DateTimeReceived: time.Now().Add(5 * time.Minute),
			OverallStatus:    "valid",
			DocumentSummary:  documents,
		})
	} else {
		ctx.JSON(400, &lhdnError{
			Code:    "BadArgument",
			Message: fmt.Sprintf("Submission ID of %s does not exist in our database", submissionUID),
		})
	}
}

func GetDocument(ctx *gin.Context) {
	documentId := ctx.Param("document_id")
	var accepted acceptedDocumentExtended
	if result := db.First(&accepted, "`uuid`", documentId); result.RowsAffected > 0 {
		ctx.JSON(200, &accepted)
	} else {
		ctx.JSON(400, &lhdnError{
			Code:    "BadArgument",
			Message: fmt.Sprintf("Document ID of %s does not exist in our database", documentId),
		})
	}
}

func UpdateDocument(ctx *gin.Context) {
	documentId := ctx.Param("document_id")
	status := ctx.Param("status")
	reason := ctx.Param("reason")

	var accepted acceptedDocumentExtended
	if result := db.First(&accepted, "`uuid`", documentId); result.RowsAffected > 0 {
		accepted.Status = status
		accepted.DocumentStatusReason = reason

		db.Save(&accepted)

		ctx.JSON(200, gin.H{
			"uuid":   accepted.UUID,
			"status": status,
		})

	} else {
		ctx.JSON(400, &lhdnError{
			Code:    "BadArgument",
			Message: fmt.Sprintf("Document ID of %s does not exist in our database", documentId),
		})
	}

}

func SubmitDocument(ctx *gin.Context) {
	contentType := ctx.GetHeader("Content-Type")
	var requestDocuments submissionFormat
	if contentType == "application/json" {
		if err := ctx.BindJSON(&requestDocuments); err != nil {
			ctx.JSON(400, &lhdnError{
				Code:    "BadArgument",
				Message: "Submission is not in JSON format",
			})
			return
		}
	} else if contentType == "application/xml" || contentType == "text/xml" {
		if err := ctx.BindXML(&requestDocuments); err != nil {
			ctx.XML(400, &lhdnError{
				Code:    "BadArgument",
				Message: "Submission is not in XML format",
			})
			return
		}
	} else {
		ctx.JSON(400, &lhdnError{
			Code:    "BadArgument",
			Message: "Unsupported format",
		})
		return
	}

	accepteds := make([]acceptedDocument, 0)
	rejecteds := make([]rejectedDocument, 0)

	submissionUID, _ := uuid.NewRandom()
	submissionUIDString := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(submissionUID[:])

	for _, doc := range requestDocuments.Documents {
		rawDocumentData, base64Err := base64.StdEncoding.DecodeString(doc.Document)
		if base64Err != nil {
			rejecteds = append(rejecteds, rejectedDocument{
				InvoiceCodeNumber: doc.CodeNumber,
				Error: lhdnError{
					Code:    "400",
					Message: "Document is not a valid Base64",
				},
			})
			continue
		}

		sha := sha256.New()
		sha.Write([]byte(rawDocumentData))
		hash := hex.EncodeToString(sha.Sum(nil))
		if hash != doc.DocumentHash {
			rejecteds = append(rejecteds, rejectedDocument{
				InvoiceCodeNumber: doc.CodeNumber,
				Error: lhdnError{
					Code:    "400",
					Message: "Document data does not match the expected Sha256 value",
				},
			})
			continue
		}

		if doc.Format == "XML" {
			var submittedDocument arbitraryXML
			if err := xml.Unmarshal(rawDocumentData, &submittedDocument); err != nil {
				rejecteds = append(rejecteds, rejectedDocument{
					InvoiceCodeNumber: doc.CodeNumber,
					Error: lhdnError{
						Code:    "400",
						Message: "Document is declared to be in XML format, but is not a valid XML",
					},
				})
				continue
			}
		} else if doc.Format == "JSON" {
			var submittedDocument json.RawMessage
			if err := json.Unmarshal(rawDocumentData, &submittedDocument); err != nil {
				fmt.Print(err)
				rejecteds = append(rejecteds, rejectedDocument{
					InvoiceCodeNumber: doc.CodeNumber,
					Error: lhdnError{
						Code:    "400",
						Message: "Document is declared to be in JSON format, but is not a valid JSON",
					},
				})
				continue
			}
		}

		var accepted acceptedDocumentExtended
		result := db.Where("`invoice_code_number` = ?", doc.CodeNumber).First(&accepted)
		if result.RowsAffected > 0 {
			rejecteds = append(rejecteds, rejectedDocument{
				InvoiceCodeNumber: doc.CodeNumber,
				Error: lhdnError{
					Code:    "422",
					Message: "Duplicate submission detected",
				},
			})
			continue
		}

		documentUuid, _ := uuid.NewRandom()
		documentUuidString := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(documentUuid[:])

		documentLongId, _ := uuid.NewRandom()
		documentLongIdString := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(documentLongId[:])

		accepteds = append(accepteds, acceptedDocument{
			UUID:              documentUuidString,
			InvoiceCodeNumber: doc.CodeNumber,
		})

		stubAccepted := acceptedDocumentExtended{
			UUID:                  documentUuidString,
			SubmissionUID:         submissionUIDString,
			LongID:                documentLongIdString,
			InvoiceCodeNumber:     doc.CodeNumber,
			TypeName:              "invoice",
			TypeVersionName:       "1.0",
			IssuerTIN:             "TAX420",
			ReceiverID:            "Banyak Untung Sdn.Bhd.",
			ReceiverName:          "Manyak Wang Sdn.Bhd.",
			DateTimeIssued:        time.Now(),
			DateTimeReceived:      time.Now(),
			DateTimeValidated:     time.Now(),
			TotalSales:            420.0,
			TotalDiscount:         69.0,
			NetAmount:             42069.0,
			Total:                 12629.0,
			Status:                "Submitted",
			CancelDateTime:        time.Now(),
			RejectRequestDateTime: time.Now(),
			DocumentStatusReason:  "Wang tarak bayar cukup",
			CreatedByUserID:       "pembekal.hebat@gmail.com",
		}

		db.Create(&stubAccepted)
	}

	output := outputDocument{
		SubmissionUID:     submissionUIDString,
		AcceptedDocuments: accepteds,
		RejectedDocuments: rejecteds,
	}

	if contentType == "application/json" {
		ctx.JSON(200, output)
	} else if contentType == "application/xml" || contentType == "text/xml" {
		ctx.XML(200, output)
	}
}
