package db

type Column int

const (
	Column_All Column = iota
	Column_Id
	Column_Description
	Column_Date
	Column_Status
	Column_ReferenceId
	Column_DateCreated
	Column_DateUpdated
	Column_ParentId
	Column_Name
	Column_Type
	Column_Currency
	Column_Active
	Column_TransactionId
	Column_AccountId
	Column_Amount
	Column_TagId
	Column_Frequency
	Column_Interval
	Column_DayOfWeek
	Column_DayOfMonth
	Column_StartDate
	Column_EndDate
	Column_Target
)

func (c Column) String() string {
	switch c {
	case Column_All:
		return "*"
	case Column_Id:
		return "id"
	case Column_Description:
		return "description"
	case Column_Date:
		return "date"
	case Column_Status:
		return "status"
	case Column_ReferenceId:
		return "reference_id"
	case Column_DateCreated:
		return "date_created"
	case Column_DateUpdated:
		return "date_updated"
	case Column_ParentId:
		return "parent_id"
	case Column_Name:
		return "name"
	case Column_Type:
		return "type"
	case Column_Currency:
		return "currency"
	case Column_Active:
		return "active"
	case Column_TransactionId:
		return "transaction_id"
	case Column_AccountId:
		return "account_id"
	case Column_Amount:
		return "amount"
	case Column_TagId:
		return "tag_id"
	case Column_Frequency:
		return "frequency"
	case Column_Interval:
		return "interval"
	case Column_DayOfWeek:
		return "day_of_week"
	case Column_DayOfMonth:
		return "day_of_month"
	case Column_StartDate:
		return "start_date"
	case Column_EndDate:
		return "end_date"
	case Column_Target:
		return "target"
	}
	return ""
}
