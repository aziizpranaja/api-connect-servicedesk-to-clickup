package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/html"
)

type CustomField struct {
	Id    string      `json:"id"`
	Value interface{} `json:"value"`
}

type ClickUpRequestBody struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Status        string        `json:"status"`
	Priority      int           `json:"priority"`
	DueDateTime   bool          `json:"due_date_time"`
	StartDateTime bool          `json:"start_date_time"`
	NotifyAll     bool          `json:"notify_all"`
	CustomFields  []CustomField `json:"custom_fields"`
}

func HitApi(c *gin.Context) {
	// var body struct {
	// 	InputData string `json:"input_data"`
	// }

	// URL tujuan
	url := "https://servicedesk.sig.id/api/v3/reports/execute_query"

	// Buat HTTP client menggunakan resty
	client := resty.New()

	// Set header untuk authentikasi menggunakan API key
	headers := map[string]string{
		"technician_Key": "6CDDD18E-3FB3-4B38-B3FA-954994F26DD2",
	}

	// Query yang kompleks
	query := `{
		"query": "SELECT wo.WORKORDERID AS 'Request ID', wo.TITLE AS 'Subject', wotodesc.FULLDESCRIPTION AS 'Description', rtdef.NAME AS 'Request Type', std.STATUSNAME AS 'Request Status', pd.PRIORITYNAME AS 'Priority', upper(ti.FIRST_NAME) AS 'Technician', CONVERT(varchar, dateadd(s, datediff(s, GETUTCDATE(), getdate()) + (wo.CREATEDTIME / 1000), '01-01-1970 07:00:00'), 20) AS 'Created Time', datename(year, dateadd(s, datediff(s, getutcdate(), getdate()) + (wo.CREATEDTIME / 1000), '1970-01-01 07:00:00')) AS 'Created Year', datename(month, dateadd(s, datediff(s, getutcdate(), getdate()) + (wo.CREATEDTIME / 1000), '1970-01-01 07:00:00')) AS 'Created Month', CONVERT(varchar, dateadd(s, datediff(s, GETUTCDATE(), getdate()) + (wo.COMPLETEDTIME / 1000), '01-01-1970 07:00:00'), 20) AS 'Completed Time', datename(year, dateadd(s, datediff(s, getutcdate(), getdate()) + (wo.COMPLETEDTIME / 1000), '1970-01-01 07:00:00')) AS 'Completed Year', datename(month, dateadd(s, datediff(s, getutcdate(), getdate()) + (wo.COMPLETEDTIME / 1000), '1970-01-01 07:00:00')) AS 'Completed Month', CONVERT(varchar, dateadd(s, datediff(s, GETUTCDATE(), getdate()) + (wo.DUEBYTIME / 1000), '01-01-1970 07:00:00'), 20) AS 'DueBy Time', sdo.NAME AS 'Site', cd.CATEGORYNAME AS 'Category', scd.NAME AS 'Subcategory', icd.NAME AS 'Item', wochangeinit.CHANGEID AS 'Change ID', wochangeinit.TITLE AS 'Change Title', stageDef.DISPLAYNAME AS 'Stage', statusDef.STATUSDISPLAYNAME AS 'Status Change', clcodeDef.NAME AS 'Change Closure Code', cmDef.FIRST_NAME AS 'Change Manager' FROM WorkOrder wo LEFT JOIN WorkOrderToDescription wotodesc ON wo.WORKORDERID =  wotodesc.WORKORDERID LEFT JOIN WorkOrderStates wos ON wo.WORKORDERID = wos.WORKORDERID LEFT JOIN SDUser td ON wos.OWNERID = td.USERID LEFT JOIN AaaUser ti ON td.USERID = ti.USER_ID LEFT JOIN StatusDefinition std ON wos.STATUSID = std.STATUSID LEFT JOIN PriorityDefinition pd ON wos.PRIORITYID = pd.PRIORITYID LEFT JOIN RequestTypeDefinition rtdef ON wos.REQUESTTYPEID = rtdef.REQUESTTYPEID LEFT JOIN CategoryDefinition cd ON wos.CATEGORYID = cd.CATEGORYID LEFT JOIN SubCategoryDefinition scd ON wos.SUBCATEGORYID = scd.SUBCATEGORYID LEFT JOIN ItemDefinition icd ON wos.ITEMID = icd.ITEMID LEFT JOIN IncidentToChangeMapping wotochangeinit ON wo.WORKORDERID = wotochangeinit.WORKORDERID LEFT JOIN SiteDefinition siteDef ON wo.SITEID = siteDef.SITEID LEFT JOIN SDOrganization sdo ON siteDef.SITEID = sdo.ORG_ID LEFT JOIN ChangeDetails wochangeinit ON wotochangeinit.CHANGEID = wochangeinit.CHANGEID LEFT JOIN Change_StageDefinition stageDef ON wochangeinit.WFSTAGEID = stageDef.WFSTAGEID LEFT JOIN Change_StatusDefinition statusDef ON wochangeinit.WFSTATUSID = statusDef.WFSTATUSID LEFT JOIN ChangeToClosureCode clcodeMapDef ON wochangeinit.CHANGEID = clcodeMapDef.CHANGEID LEFT JOIN Change_ClosureCode clcodeDef ON clcodeMapDef.ID = clcodeDef.ID LEFT JOIN AaaUser cmDef ON wochangeinit.CHANGEMANAGERID = cmDef.USER_ID WHERE rtdef.NAME = 'Demand Request' AND sdo.NAME NOT IN ('Solusi Bangun Indonesia.', 'PT. Semen Baturaja', 'PT Sinergi Informatika Semen Indonesia.') AND icd.NAME NOT IN ('Infrasructure','Support Infrasructure','Infrastructure','Desktop/PC','Laptop','Service Desk Plus (Manage Engine)','Transport','Update/Edit Data','User Guide','Forca Employee Self Service','Success Factor') AND scd.NAME NOT IN ('Hardware Services','Permintaan Software') AND cd.CATEGORYNAME NOT IN ('IT - Infrastructure Services','IT - Account Management','IT - Security','01. Aplikasi SMBR') AND ((wo.CREATEDTIME >= datetolong('2021-01-01') AND wo.CREATEDTIME IS NOT NULL AND wo.CREATEDTIME != 0) AND (wo.CREATEDTIME <= datetolong('2023-12-31') AND wo.CREATEDTIME IS NOT NULL AND wo.CREATEDTIME != 0 AND wo.CREATEDTIME != -1)) AND ((wo.SITEID IN (SELECT UserSiteMapping.SITEID FROM UserSiteMapping WHERE (UserSiteMapping.USERID = 97695))) OR ((wos.OWNERID = 97695) OR (wo.REQUESTERID = 97695))) AND wo.ISPARENT = '1' AND wo.TITLE != 'Permintaan Reset Device Forca ESS' AND stageDef.DISPLAYNAME != 'Close' AND clcodeDef.NAME IS NULL"
	  }`

	// Kirim permintaan POST ke API dengan form-data
	resp, err := client.R().
		SetHeaders(headers).
		SetFormData(map[string]string{
			// "input_data": body.InputData,
			"input_data": query,
		}).
		Post(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ubah respons menjadi bentuk JSON
	var responseData map[string]interface{}
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	// Ambil bagian "data" dari respons
	data := responseData["execute_query"].(map[string]interface{})["data"]

	// url untuk hit clickup
	urlClickup := "https://api.clickup.com/api/v2/list/901800049325/task"

	// Set header untuk autentikasi menggunakan ClickUp
	headers_clickup := map[string]string{
		"Authorization": "Bearer 66807434_a7e9bfdb71eae6260fd2ccfe60157956f3864daacc7d00a08356f9390b915a4e",
	}

	// Body permintaan
	requestBody := ClickUpRequestBody{
		Name:          "percobaan buat task 1",
		Description:   "percobaan buat task 1",
		Status:        "TO DO",
		Priority:      3,
		DueDateTime:   false,
		StartDateTime: false,
		NotifyAll:     true,
		CustomFields: []CustomField{
			{
				Id:    "0a52c486-5f05-403b-b4fd-c512ff05131c",
				Value: "test",
			},
			{
				Id:    "03efda77-c7a0-42d3-8afd-fd546353c2f5",
				Value: "Text field input",
			},
		},
	}

	// Loop melalui setiap elemen data untuk diinputkan kedalam body
	for _, item := range data.([]interface{}) {
		// Konversi item menjadi map[string]interface{}
		taskData := item.(map[string]interface{})

		// complete time
		var completedTime, _ = time.Parse("2006-01-02 15:04:05", taskData["Completed Time"].(string))
		unixTimeMillis := completedTime.UnixNano() / int64(time.Millisecond)

		// created time
		var createdTime, _ = time.Parse("2006-01-02 15:04:05", taskData["Created Time"].(string))
		unixTimeMillisCreated := createdTime.UnixNano() / int64(time.Millisecond)

		// dueby time
		var duebyTime, _ = time.Parse("2006-01-02 15:04:05", taskData["DueBy Time"].(string))
		unixTimeMillisDue := duebyTime.UnixNano() / int64(time.Millisecond)

		// Request Status
		var requestStatus = ""
		if taskData["Request Status"].(string) == "Closed" {
			requestStatus = "5c37c39c-e5d9-4441-9481-bb7c78a094cb"
		} else if taskData["Request Status"].(string) == "Waiting Approval ICT" {
			requestStatus = "017e1da3-cacc-453b-8f9e-5e8e934d09f2"
		} else if taskData["Request Status"].(string) == "Canceled" {
			requestStatus = "edad8dda-b98a-4ca1-ac40-26cf7b2d8104"
		} else if taskData["Request Status"].(string) == "Onhold" {
			requestStatus = "fc543273-da3c-4de3-b886-c6702edbbaea"
		} else if taskData["Request Status"].(string) == "In Progress - Change Request" {
			requestStatus = "d2b31c17-a80f-4490-9738-070704496a45"
		} else if taskData["Request Status"].(string) == "In Progress" {
			requestStatus = "72fb4ccd-0bd5-4e4e-bd36-d225f76566fc"
		} else if taskData["Request Status"].(string) == "In Progress - Waiting End User Confirmation" {
			requestStatus = "6261811f-33c4-495c-b304-59f26e490e11"
		} else if taskData["Request Status"].(string) == "In Progress - Project" {
			requestStatus = "d910661d-d7e2-4095-8c39-5c7d3cbeb041"
		} else if taskData["Request Status"].(string) == "Resolved" {
			requestStatus = "cc62f939-c358-4497-8936-f62b3183c1e0"
		} else if taskData["Request Status"].(string) == "In Progress - On Evaluation" {
			requestStatus = "b1cb0595-fe39-4892-8987-49f6baa62c32"
		} else if taskData["Request Status"].(string) == "Waiting Approval" {
			requestStatus = "a9b4c373-5566-4ff3-be30-7e4768d3109d"
		} else if taskData["Request Status"].(string) == "Open" {
			requestStatus = "bb13477e-758b-4403-8b31-5aa8f3cb78b0"
		}

		// site
		var site = ""
		if taskData["Site"].(string) == "PT. Semen Indonesia." {
			site = "ce570d21-f49c-4433-8a18-18495617ca5c"
		} else if taskData["Site"].(string) == "PT. Semen Padang." {
			site = "e94b789f-1999-464c-83a3-580d3e0847c8"
		} else if taskData["Site"].(string) == "PT. Semen Tonasa." {
			site = "0ee47888-bda5-4ef9-8343-2e9c644d7774"
		} else if taskData["Site"].(string) == "PT. Semen Gresik." {
			site = "22c36640-63d3-4430-9f3f-a5219b19402c"
		} else if taskData["Site"].(string) == "Thang Long Cement Company" {
			site = "e27ec1fa-9cf3-412c-b5b2-1c48406e5863"
		}

		// Priorities
		var priority = ""
		if taskData["Priority"].(string) == "P1 - Critical" {
			priority = "ef0bd234-a8ff-4633-b315-70806359c5ee"
		} else if taskData["Priority"].(string) == "P2 - High" {
			priority = "158a7d94-8f9a-417d-9217-6e3ad1e17f93"
		} else if taskData["Priority"].(string) == "P3 - Medium" {
			priority = "94f3c545-3946-43be-b865-bd19b17d8633"
		} else if taskData["Priority"].(string) == "P4 - Normal" {
			priority = "1add6698-41d4-4671-b734-657bfc6ce013"
		} else if taskData["Priority"].(string) == "P5 - Low" {
			priority = "4a2db263-c82e-490c-98c7-b819a8ed111c"
		}

		// Status Change
		var statusChange = ""
		if taskData["Status Change"].(string) == "Completed" {
			statusChange = "145588dc-de26-49cb-8502-a00c1b4f368d"
		} else if taskData["Status Change"].(string) == "Canceled" {
			statusChange = "14937c59-1875-4e3b-936d-623303c26a90"
		} else if taskData["Status Change"].(string) == "Approval Pending" {
			statusChange = "91350d40-b7e5-419d-ac30-d487f2d1b845"
		} else if taskData["Status Change"].(string) == "In Progress" {
			statusChange = "cbebdc43-c757-4a6f-9b74-d0977e42022c"
		} else if taskData["Status Change"].(string) == "On Hold" {
			statusChange = "e7bb06c8-f0dd-4e2c-b3f1-84bf131a574d"
		} else if taskData["Status Change"].(string) == "Planning In Progress" {
			statusChange = "768eb73a-b51d-4611-89eb-3feb692e732e"
		} else if taskData["Status Change"].(string) == "Requested" {
			statusChange = "08819019-4b3f-453a-bb5e-3027f4ab26ab"
		} else if taskData["Status Change"].(string) == "Rejected" {
			statusChange = "27d22edb-6227-441d-ae63-7e460d42162c"
		}

		// Change Manager
		var changeManager = ""
		if taskData["Change Manager"].(string) == "ICT SIG - Nanang Iqbal Habibie" {
			changeManager = "6bcad070-aa37-41e4-a68b-211bfe451821"
		} else if taskData["Change Manager"].(string) == "ICT SIG - M. ZAINUL A." {
			changeManager = "f98ab4d5-3c7e-45f2-8c0b-0b87541ccb24"
		} else if taskData["Change Manager"].(string) == "ICT SIG - Evy Dahniar" {
			changeManager = "b61704d5-3231-4d56-b42f-59765f0b15cb"
		} else if taskData["Change Manager"].(string) == "ICT SIG - AHMAD BAIHAQI" {
			changeManager = "9f6688db-ccca-484e-af02-5c3e2b1041be"
		} else if taskData["Change Manager"].(string) == "ICT SIG - RONALIS A." {
			changeManager = "77e684ca-3d2a-4bd4-b155-6481abb028e4"
		} else if taskData["Change Manager"].(string) == "ICT SIG - BAMBANG TRIMONO" {
			changeManager = "a923bb72-ec6e-4053-8377-77a993fc2e12"
		} else if taskData["Change Manager"].(string) == "ICT SIG - YULMIZAR" {
			changeManager = "7fd27278-1e20-4d1a-942f-8946a7cda99f"
		} else if taskData["Change Manager"].(string) == "ICT SIG - Eko Dharmawan Nizar" {
			changeManager = "87233164-b20d-4cb2-8063-94011210286c"
		}

		// created year
		var createdYear = ""
		if taskData["Created Year"].(string) == "2021" {
			createdYear = "d865ac43-b5dc-4c8f-a805-0abbba02d4ab"
		} else if taskData["Created Year"].(string) == "2022" {
			createdYear = "543effeb-cda2-49c4-a86a-a05b4de8cf24"
		} else if taskData["Created Year"].(string) == "2023" {
			createdYear = "cac1a4db-dd63-46db-8726-ca1c1164b36a"
		} else if taskData["Created Year"].(string) == "2024" {
			createdYear = "dda1dfea-db66-488b-9020-20caa8ef953c"
		} else if taskData["Created Year"].(string) == "2025" {
			createdYear = "61df785a-20f3-4ee9-bbdc-af0693eca522"
		}

		// complete year
		var completeYear = ""
		if taskData["Completed Year"].(string) == "2021" {
			completeYear = "8df8992a-1eed-45a6-9bd4-58e8dc87b2b5"
		} else if taskData["Completed Year"].(string) == "2022" {
			completeYear = "035dbaa6-3049-44d6-bbe3-bfbaf755236f"
		} else if taskData["Completed Year"].(string) == "2023" {
			completeYear = "ce8c5ad0-6c4e-401b-a701-58a771c434c5"
		} else if taskData["Completed Year"].(string) == "2024" {
			completeYear = "75657002-2a69-4519-8a8d-7a31316f45cd"
		} else if taskData["Completed Year"].(string) == "2025" {
			completeYear = "73fbbeb9-7929-4df0-8e38-19e3e57d4c1b"
		} else {
			completeYear = ""
		}

		// created month
		var createdMonth = ""
		if taskData["Created Month"].(string) == "January" {
			createdMonth = "866725cb-4d24-435a-b080-0dea0f66e41c"
		} else if taskData["Created Month"].(string) == "February" {
			createdMonth = "fe1dad5f-634d-4179-9030-f269ce768479"
		} else if taskData["Created Month"].(string) == "March" {
			createdMonth = "f9c3a4e2-43d6-467c-a77e-63c63b81b178"
		} else if taskData["Created Month"].(string) == "April" {
			createdMonth = "c8459015-ac28-4bec-b824-08cffbd3d6e3"
		} else if taskData["Created Month"].(string) == "May" {
			createdMonth = "8d64db82-4bd2-4b83-be69-dfb026d8938d"
		} else if taskData["Created Month"].(string) == "June" {
			createdMonth = "a8c08369-f171-4ebe-b444-bdf4adc60744"
		} else if taskData["Created Month"].(string) == "July" {
			createdMonth = "0d7e799e-ce10-42b3-b9cb-9f72b8a63eb8"
		} else if taskData["Created Month"].(string) == "August" {
			createdMonth = "293d3a06-a7eb-4fde-9a2f-733fd5c101da"
		} else if taskData["Created Month"].(string) == "September" {
			createdMonth = "af856110-e79a-496f-b735-c47ba1eb9939"
		} else if taskData["Created Month"].(string) == "October" {
			createdMonth = "65ade436-66c2-41aa-a6c2-f0dd5f1f16ef"
		} else if taskData["Created Month"].(string) == "November" {
			createdMonth = "8670f77d-e2c4-43e4-9957-196c2f61f5bb"
		} else if taskData["Created Month"].(string) == "December" {
			createdMonth = "d33b42ea-bb48-40bc-954e-80a53ef86516"
		}

		// complete month
		var completeMonth = ""
		if taskData["Completed Month"].(string) == "January" {
			completeMonth = "19d220fb-0baa-4d1c-8c49-5edd397d16d2"
		} else if taskData["Completed Month"].(string) == "February" {
			completeMonth = "d5d212ac-f5e1-41dd-87d8-108690e7a30a"
		} else if taskData["Completed Month"].(string) == "March" {
			completeMonth = "9275bc25-bb3d-4d2d-9887-85147908f55d"
		} else if taskData["Completed Month"].(string) == "April" {
			completeMonth = "8bd7818f-3b29-4967-b6fa-b73eb7fd347a"
		} else if taskData["Completed Month"].(string) == "May" {
			completeMonth = "53811cfe-5686-4ed1-acc6-b0c413c1c064"
		} else if taskData["Completed Month"].(string) == "June" {
			completeMonth = "ff00b08a-5be2-4a26-95e6-b988e4003814"
		} else if taskData["Completed Month"].(string) == "July" {
			completeMonth = "4a387796-5dab-4826-bbb7-c6d126d3789f"
		} else if taskData["Completed Month"].(string) == "August" {
			completeMonth = "c6a438e6-a0b0-42d4-9d92-e97b7aa24498"
		} else if taskData["Completed Month"].(string) == "September" {
			completeMonth = "35155e55-9f88-4939-a037-ba90b7b3d1c0"
		} else if taskData["Completed Month"].(string) == "October" {
			completeMonth = "7ade2b0b-3891-4f04-9fe0-9b716f916277"
		} else if taskData["Completed Month"].(string) == "November" {
			completeMonth = "65d5ba05-b35e-42f1-bde5-a1c26ae260f0"
		} else if taskData["Completed Month"].(string) == "December" {
			completeMonth = "0b4ba3f1-4cea-4bef-aa34-e638bba5c593"
		}

		// stage
		var stage = ""
		if taskData["Stage"].(string) == "Submission" {
			stage = "b551f340-a782-4a47-a0f2-9df2771a0eb2"
		} else if taskData["Stage"].(string) == "Planning" {
			stage = "e9f1b20f-b103-4ab9-9204-b6db41f3f368"
		} else if taskData["Stage"].(string) == "Approval" {
			stage = "d9644f2f-7b0c-4e4a-a980-3ffd797af6bd"
		} else if taskData["Stage"].(string) == "Implementation" {
			stage = "c243aa4d-9a2d-483f-bd0f-fbd27bcccac7"
		} else if taskData["Stage"].(string) == "Review" {
			stage = "88487237-9ef4-4803-be4f-52a277c1f73c"
		} else if taskData["Stage"].(string) == "Close" {
			stage = "e2ba8e02-ed1f-45cd-a88b-f8d2ce45f525"
		} else {
			stage = ""
		}
		// htmlContent := taskData["Description"].(string)

		// translatedText := TranslateHTML(htmlContent)

		// Ubah nilai body permintaan sesuai dengan data saat ini
		requestBody.Name = taskData["Subject"].(string)
		requestBody.Description = ""
		requestBody.Status = "TO DO"
		requestBody.Priority = 3
		requestBody.DueDateTime = false
		requestBody.StartDateTime = false
		requestBody.NotifyAll = true

		// Ganti nilai CustomField
		requestBody.CustomFields = []CustomField{
			{
				// Request ID
				Id:    "7650f3a9-a3eb-43f1-b66c-8f25e9de16dd",
				Value: fmt.Sprintf("%v", taskData["Request ID"]),
			},
			{
				// request type
				Id:    "31db95a8-4a03-4794-81df-5d9d8ceb3ffa",
				Value: taskData["Request Type"].(string),
			},
			{
				// Request status
				Id:    "f03c3c49-d8f4-4a6a-8300-afcff6b4fb81",
				Value: requestStatus,
			},
			{
				// priority
				Id:    "635a28b9-7953-494e-8f20-6d81345eb33d",
				Value: priority,
			},
			{
				// technician
				Id:    "861763ef-a215-4573-996b-5f62f55dff4e",
				Value: taskData["Technician"].(string),
			},
			{
				// created time
				Id:    "fa5272a3-bb98-4ea6-8571-bcdba9b26194",
				Value: unixTimeMillisCreated,
			},
			{
				// created year
				Id:    "6297f187-63a4-4457-b783-1954abbbbcda",
				Value: createdYear,
			},
			{
				// created month
				Id:    "7546d8ab-5797-4bcd-bd28-980feec23b7a",
				Value: createdMonth,
			},
			{
				// completed time
				Id:    "f416365c-7aed-4132-93f2-3e3c93a3446a",
				Value: unixTimeMillis,
			},
			{
				// completed year
				Id:    "39bff84f-eef9-447e-8446-9c27299ac23a",
				Value: completeYear,
			},
			{
				// complete month
				Id:    "f2e868c5-1a71-4536-a594-3e6bd839052b",
				Value: completeMonth,
			},
			{
				// dueby time
				Id:    "8b78b274-292e-4282-907a-c359e50d099c",
				Value: unixTimeMillisDue,
			},
			{
				// site
				Id:    "f5794a0a-e530-42c7-89bd-f7951fe44d2c",
				Value: site,
			},
			{
				// category
				Id:    "b67663d1-bc46-49f4-b5e0-51fb7f7cc6c7",
				Value: taskData["Category"].(string),
			},
			{
				// subcategory
				Id:    "8111d57d-206f-479f-88fe-188860ab19a6",
				Value: taskData["Subcategory"].(string),
			},
			{
				// item
				Id:    "8a2201b2-ce07-4989-baf2-66297c4cc43c",
				Value: taskData["Item"].(string),
			},
			{
				// change id
				Id:    "b7f964ef-50c5-462c-8ad0-3d439effbd38",
				Value: fmt.Sprintf("%v", taskData["Change ID"]),
			},
			{
				// status change
				Id:    "771f8134-d77b-4b08-906a-af2d69f5631e",
				Value: statusChange,
			},
			{
				// change manager
				Id:    "eb087393-09f4-4750-a436-af2796e8de7f",
				Value: changeManager,
			},
			// {
			// 	// direktorat
			// 	Id:    "65cdc449-5819-4a07-a4dc-0efa1a5fd02d",
			// 	Value: "-",
			// },
			{
				// change title
				Id:    "baf04c06-fd4c-4540-b5ca-3584390a77c7",
				Value: taskData["Change Title"].(string),
			},
			{
				// change closure
				Id:    "3f65c693-5963-4c88-ab81-9c5870624257",
				Value: "-",
			},
			{
				// stage
				Id:    "b861ac13-97ce-4be2-bb0b-d4a88ea0c330",
				Value: stage,
			},
		}

		// Konversi body menjadi JSON
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
			return
		}

		// Kirim permintaan POST dengan body JSON dan header Authorization
		_, err = client.R().
			SetHeaders(headers_clickup).
			SetHeader("Content-Type", "application/json").
			SetBody(jsonBody).
			Post(urlClickup)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Run Function",
	})
}

func TranslateHTML(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		// Penanganan kesalahan jika terjadi
	}

	var translatedText string
	var translateNode func(*html.Node)
	translateNode = func(n *html.Node) {
		if n.Type == html.TextNode {
			translatedText += n.Data
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			translateNode(c)
		}
	}

	translateNode(doc)

	return translatedText
}
