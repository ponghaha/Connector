package format

import (
	"strconv"

	"connectorapi-go/internal/core/domain"
	"connectorapi-go/internal/adapter/utils"
)

// Converts SubmitLoanApplicationRequest to a fixed-length string.
func FormatSubmitLoanApplicationRequest(submitLoanApplicationReq domain.SubmitLoanApplicationRequest) string {
	requestID                      := utils.PadOrTruncate(submitLoanApplicationReq.RequestID, 20)
	applicationNo                  := utils.PadOrTruncate(submitLoanApplicationReq.ApplicationNo, 20)
	keptBoxNo                      := utils.PadOrTruncate(submitLoanApplicationReq.KeptBoxNo, 15)
	customerGroup                  := utils.PadOrTruncate(submitLoanApplicationReq.CustomerGroup, 1)
	nonMemberType                  := utils.PadOrTruncate(submitLoanApplicationReq.NonmemberType, 1)
	iamsApplicationReceivedDateTime:= utils.PadOrTruncate(submitLoanApplicationReq.ApplicationDate, 14)
	ncbToken                       := utils.PadOrTruncate(submitLoanApplicationReq.NCBToken, 25)
	customerID                     := utils.PadOrTruncate(submitLoanApplicationReq.IDCardNo, 20)
	customerEngPrefix              := utils.PadOrTruncate(submitLoanApplicationReq.TitleNameEN, 20)
	customerEngName                := utils.PadOrTruncate(submitLoanApplicationReq.NameEN, 30)
	customerThaiName               := utils.PadOrTruncate(submitLoanApplicationReq.NameTH, 30)
	age                            := utils.PadIntWithZero(submitLoanApplicationReq.Age, 3)
	birthdate                      := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.Birthdate), 8)
	gender                         := utils.PadOrTruncate(strconv.Itoa(submitLoanApplicationReq.Gender), 1)
	marriage                       := utils.PadOrTruncate(strconv.Itoa(submitLoanApplicationReq.MarriageStatus), 1)
	homeNumber                     := utils.PadOrTruncate(submitLoanApplicationReq.HomeAddressNo, 20)
	homeMoo                        := utils.PadOrTruncate(submitLoanApplicationReq.HomeMoo, 3)
	homeVillageName                := utils.PadOrTruncate(submitLoanApplicationReq.HomeVillage, 45)
	homeRoom                       := utils.PadOrTruncate(submitLoanApplicationReq.HomeRoom, 10)
	homeFloor                      := utils.PadOrTruncate(submitLoanApplicationReq.HomeFloor, 10)
	homeSoi                        := utils.PadOrTruncate(submitLoanApplicationReq.HomeSoi, 25)
	homeRoad                       := utils.PadOrTruncate(submitLoanApplicationReq.HomeRoad, 25)
	homeSuburb                     := utils.PadOrTruncate(submitLoanApplicationReq.HomeSubDistrict, 25)
	homeDistrict                   := utils.PadOrTruncate(submitLoanApplicationReq.HomeDistrict, 25)
	homeProvince                   := utils.PadOrTruncate(submitLoanApplicationReq.HomeProvince, 20)
	homeZip                        := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.HomeZipCode), 5)
	homePhone                      := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.HomePhone), 10)
	homePhoneExtension             := utils.PadOrTruncate(submitLoanApplicationReq.HomePhoneExt, 5)
	mobileNumber                   := utils.PadOrTruncate(submitLoanApplicationReq.MobileNo, 15)
	homeStatus                     := utils.PadOrTruncate(submitLoanApplicationReq.HomeStatus, 1)
	emailAddress                   := utils.PadOrTruncate(submitLoanApplicationReq.Email, 35)
	livingPeriod                   := utils.PadFloatWithZero(submitLoanApplicationReq.LivingPeriod, 4, 2)
	stayWith                       := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.StayWith), 3)
	educationCode                  := utils.PadOrTruncate(submitLoanApplicationReq.EducationCode, 2)
	educationDescription           := utils.PadOrTruncate(submitLoanApplicationReq.EducationDescription, 20)
	officeName                     := utils.PadOrTruncate(submitLoanApplicationReq.OfficeName, 50)
	officeSection                  := utils.PadOrTruncate(submitLoanApplicationReq.OfficeSection, 30)
	officeNumber                   := utils.PadOrTruncate(submitLoanApplicationReq.OfficeNo, 20)
	officeMoo                      := utils.PadOrTruncate(submitLoanApplicationReq.OfficeMoo, 2)
	officeBuildingName             := utils.PadOrTruncate(submitLoanApplicationReq.OfficeBuildingName, 45)
	officeRoom                     := utils.PadOrTruncate(submitLoanApplicationReq.OfficeRoom, 10)
	officeFloorNumber              := utils.PadOrTruncate(submitLoanApplicationReq.OfficeFloor, 10)
	officeSoi                      := utils.PadOrTruncate(submitLoanApplicationReq.OfficeSoi, 25)
	officeRoad                     := utils.PadOrTruncate(submitLoanApplicationReq.OfficeRoad, 25)
	officeSuburb                   := utils.PadOrTruncate(submitLoanApplicationReq.OfficeSubDistrict, 25)
	officeDistrict                 := utils.PadOrTruncate(submitLoanApplicationReq.OfficeDistrict, 25)
	officeProvince                 := utils.PadOrTruncate(submitLoanApplicationReq.OfficeProvince, 20)
	officeZip                      := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.OfficeZipCode), 5)
	officePhone                    := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.OfficePhone), 10)
	officeExtension                := utils.PadOrTruncate(submitLoanApplicationReq.OfficePhoneExt, 5)
	jobType                        := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.JobTypeCode), 2)
	jobTypeDetail                  := utils.PadOrTruncate(submitLoanApplicationReq.JobTypeDetailCode, 2)
	workingPeriod                  := utils.PadFloatWithZero(submitLoanApplicationReq.WorkingPeriod, 4, 2)
	employmentStatus               := utils.PadOrTruncate(submitLoanApplicationReq.EmploymentStatus, 2)
	businessType                   := utils.PadOrTruncate(submitLoanApplicationReq.BusinessType, 2)
	businessDescription            := utils.PadOrTruncate(submitLoanApplicationReq.BusinessDescription, 20)
	timeToContract                 := utils.PadOrTruncate(submitLoanApplicationReq.TimeToContact, 30)
	spouseName                     := utils.PadOrTruncate(submitLoanApplicationReq.SpouseName, 30)
	spousePhone                    := utils.PadOrTruncate(submitLoanApplicationReq.SpousePhone, 10)
	spouseExtension                := utils.PadOrTruncate(submitLoanApplicationReq.SpousePhoneExt, 5)
	debtReferenceName              := utils.PadOrTruncate(submitLoanApplicationReq.DebtReferenceName, 30)
	debtReferencerelationship      := utils.PadOrTruncate(submitLoanApplicationReq.DebtReferenceRelationship, 15)
	debtReferencePhone             := utils.PadOrTruncate(submitLoanApplicationReq.DebtReferencePhone, 11)
	debtReferencePhoneExtension    := utils.PadOrTruncate(submitLoanApplicationReq.DebtReferencePhoneExt, 5)
	debtReferenceMobilePhone       := utils.PadOrTruncate(submitLoanApplicationReq.DebtReferenceMobile, 11)
	salary                         := utils.PadFloatWithZero(submitLoanApplicationReq.Salary, 11, 2)
	otherIncome                    := utils.PadFloatWithZero(submitLoanApplicationReq.OtherIncome, 11, 2)
	otherIncomeResource            := utils.PadOrTruncate(submitLoanApplicationReq.OtherIncomResource, 1)
	otherIncomeResourceDescription := utils.PadOrTruncate(submitLoanApplicationReq.OtherincomResourceDesc, 20)
	paymentType                    := utils.PadOrTruncate(submitLoanApplicationReq.PaymentType, 2)
	autoPayBank                    := utils.PadOrTruncate(submitLoanApplicationReq.AutopayBank, 30)
	autoPayAccountNo               := utils.PadOrTruncate(submitLoanApplicationReq.AutopayAccountNo, 10)
	mailTo                         := utils.PadOrTruncate(submitLoanApplicationReq.MailTo, 1)
	idIssuedDate                   := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.IDCardIssueDate), 8)
	idExpiryDate                   := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.IDCardExpiryDate), 8)
	idHouseNumber                  := utils.PadOrTruncate(submitLoanApplicationReq.IDCardHouseNo, 20)
	idMoo                          := utils.PadOrTruncate(submitLoanApplicationReq.IDCardMoo, 2)
	idRoom                         := utils.PadOrTruncate(submitLoanApplicationReq.IDCardRoom, 10)
	idFloor                        := utils.PadOrTruncate(submitLoanApplicationReq.IDCardFloor, 10)
	idSoi                          := utils.PadOrTruncate(submitLoanApplicationReq.IDCardSoi, 25)
	idRoad                         := utils.PadOrTruncate(submitLoanApplicationReq.IDCardRoad, 25)
	idSuburb                       := utils.PadOrTruncate(submitLoanApplicationReq.IDCardSubDistrict, 25)
	idDistrict                     := utils.PadOrTruncate(submitLoanApplicationReq.IDCardDistrict, 25)
	idProvince                     := utils.PadOrTruncate(submitLoanApplicationReq.IDCardProvince, 20)
	agentCode                      := utils.PadOrTruncate(submitLoanApplicationReq.AgentCode, 8)
	applicationPurpose             := utils.PadOrTruncate(submitLoanApplicationReq.ApplicationPurposeCode, 3)
	otherPurposeDescription        := utils.PadOrTruncate(submitLoanApplicationReq.OtherPurposeDescription, 100)
	applyType                      := utils.PadOrTruncate(submitLoanApplicationReq.ApplyType, 3)
	productCode                    := utils.PadOrTruncate(submitLoanApplicationReq.ProductCode, 4)
	makerCode                      := utils.PadOrTruncate(submitLoanApplicationReq.BrandCode, 4)
	modelCode                      := utils.PadOrTruncate(submitLoanApplicationReq.ModelCode, 10)
	color                          := utils.PadOrTruncate(submitLoanApplicationReq.Color, 20)
	cc                             := utils.PadIntWithZero(submitLoanApplicationReq.Cc, 5)
	powerTransitionType            := utils.PadOrTruncate(submitLoanApplicationReq.PowerTransitionType, 1)
	newOrUsedCard                  := utils.PadOrTruncate(submitLoanApplicationReq.CarusedType, 1)
	carMileageNo                   := utils.PadIntWithZero(submitLoanApplicationReq.CarmileNo, 10)
	enginePlateNo                  := utils.PadOrTruncate(submitLoanApplicationReq.EngineNo, 30)
	chassisPlateNo                 := utils.PadOrTruncate(submitLoanApplicationReq.ChassisNo, 20)
	registrationDate               := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.CarRegistrationDate), 8)
	registrationNo                 := utils.PadOrTruncate(submitLoanApplicationReq.LicensePlateNo, 30)
	registrationProvince           := utils.PadOrTruncate(submitLoanApplicationReq.LicensePlateProvince, 3)
	hasLoanContract                := utils.PadOrTruncate(submitLoanApplicationReq.LoanContractFlag, 1)
	estimationPrice                := utils.PadFloatWithZero(submitLoanApplicationReq.EstimationPrice, 11, 2)
	cashPrice                      := utils.PadFloatWithZero(submitLoanApplicationReq.CashPrice, 11, 2)
	promotion                      := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.PromotionCode), 10)
	rate                           := utils.PadFloatWithZero(submitLoanApplicationReq.InterestRate, 5, 3)
	down                           := utils.PadOrTruncate(submitLoanApplicationReq.DownPayment, 1)
	financePrice                   := utils.PadFloatWithZero(submitLoanApplicationReq.FinancePrice, 11, 2)
	downPayment                    := utils.PadFloatWithZero(submitLoanApplicationReq.DownPaymentPrice, 10, 2)
	term                           := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.InstallmentPeriod), 3)
	carType                        := utils.PadOrTruncate(submitLoanApplicationReq.CarType, 2)
	note                           := utils.PadOrTruncate(submitLoanApplicationReq.Note, 50)
	scannedBy                      := utils.PadOrTruncate(submitLoanApplicationReq.ScanedBy, 30)
	marketingCode                  := utils.PadOrTruncate(submitLoanApplicationReq.MarketingCode, 30)
	manufactureYear                := utils.PadIntWithZero(utils.ConvertStringToInt(submitLoanApplicationReq.ManufactureYear), 4)
	applicationReceivedDateTime    := utils.PadOrTruncate(submitLoanApplicationReq.ApplicationReceivedDate, 14)
	bankCode                       := utils.PadOrTruncate(submitLoanApplicationReq.BankCode, 3)
	accountNo                      := utils.PadOrTruncate(submitLoanApplicationReq.AccountNo, 20)
	acHolderNameThai               := utils.PadOrTruncate(submitLoanApplicationReq.AccountHolderName, 30)

	return 	requestID + applicationNo + keptBoxNo + customerGroup + nonMemberType +
			iamsApplicationReceivedDateTime + ncbToken + customerID + customerEngPrefix +customerEngName +
			customerThaiName + age + birthdate + gender + marriage +
			homeNumber + homeMoo + homeVillageName + homeRoom + homeFloor +
			homeSoi + homeRoad + homeSuburb + homeDistrict + homeProvince +
			homeZip + homePhone + homePhoneExtension + mobileNumber + homeStatus +
			emailAddress + livingPeriod + stayWith + educationCode + educationDescription +
			officeName + officeSection + officeNumber + officeMoo + officeBuildingName +
			officeRoom + officeFloorNumber + officeSoi + officeRoad + officeSuburb +
			officeDistrict + officeProvince + officeZip + officePhone + officeExtension +
			jobType + jobTypeDetail + workingPeriod + employmentStatus + businessType +
			businessDescription + timeToContract + spouseName + spousePhone + spouseExtension +
			debtReferenceName + debtReferencerelationship + debtReferencePhone + debtReferencePhoneExtension + debtReferenceMobilePhone +
			salary + otherIncome + otherIncomeResource + otherIncomeResourceDescription + paymentType +
			autoPayBank + autoPayAccountNo + mailTo + idIssuedDate + idExpiryDate +
			idHouseNumber + idMoo + idRoom + idFloor + idSoi +
			idRoad + idSuburb + idDistrict + idProvince + agentCode +
			applicationPurpose + otherPurposeDescription + applyType + productCode + makerCode +
			modelCode + color + cc + powerTransitionType + newOrUsedCard +
			carMileageNo + enginePlateNo + chassisPlateNo + registrationDate + registrationNo +
			registrationProvince + hasLoanContract + estimationPrice + cashPrice + promotion +
			rate + down + financePrice + downPayment + term +
			carType + note + scannedBy + marketingCode + manufactureYear +
			applicationReceivedDateTime + bankCode + accountNo + acHolderNameThai
}
