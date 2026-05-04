package system

import "erp-service/common"

// Contract 合同主表
type Contract struct {
	ID           uint             `json:"id" gorm:"primaryKey"`
	ProjectName  string           `json:"projectName" gorm:"column:project_name;size:255"`
	OrderNo      string           `json:"orderNo" gorm:"column:order_no;size:100"`
	OrderDate    string           `json:"orderDate" gorm:"column:order_date;size:20"`
	FromCompany  string           `json:"fromCompany" gorm:"column:from_company;size:200"`
	ToCompany    string           `json:"toCompany" gorm:"column:to_company;size:200"`
	Buyer        string           `json:"buyer" gorm:"column:buyer;size:50"`
	Attn         string           `json:"attn" gorm:"column:attn;size:50"`
	BuyerEmail   string           `json:"buyerEmail" gorm:"column:buyer_email;size:100"`
	BuyerTel     string           `json:"buyerTel" gorm:"column:buyer_tel;size:50"`
	AttnTel      string           `json:"attnTel" gorm:"column:attn_tel;size:50"`
	TotalAmount  string           `json:"totalAmount" gorm:"column:total_amount;size:30"`
	DeliveryAddr string           `json:"deliveryAddr" gorm:"column:delivery_addr;size:500"`
	Remark       string           `json:"remark" gorm:"column:remark;size:500"`
	CreatedBy    uint             `json:"createdBy" gorm:"column:created_by"`
	CreatedAt    common.MilliTime `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    common.MilliTime `json:"updatedAt" gorm:"column:updated_at;->"`
}

func (Contract) TableName() string { return common.TableContract }

// ContractItem 合同明细
type ContractItem struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	ContractID uint   `json:"contractId" gorm:"column:contract_id;index"`
	Seq        int    `json:"seq" gorm:"column:seq"`
	Name       string `json:"name" gorm:"column:name;size:200"`
	Spec       string `json:"spec" gorm:"column:spec;size:200"`
	Brand      string `json:"brand" gorm:"column:brand;size:100"`
	Qty        string `json:"qty" gorm:"column:qty;size:30"`
	Unit       string `json:"unit" gorm:"column:unit;size:20"`
	UnitPrice  string `json:"unitPrice" gorm:"column:unit_price;size:30"`
	Amount     string `json:"amount" gorm:"column:amount;size:30"`
	Operator   string `json:"operator" gorm:"column:operator;size:50"`
	Location   string `json:"location" gorm:"column:location;size:100"`
	Remark     string `json:"remark" gorm:"column:remark;size:200"`
}

func (ContractItem) TableName() string { return common.TableContractItem }

// --- DTO ---

type ContractItemReq struct {
	Seq       int    `json:"seq"`
	Name      string `json:"name"`
	Spec      string `json:"spec"`
	Brand     string `json:"brand"`
	Qty       string `json:"qty"`
	Unit      string `json:"unit"`
	UnitPrice string `json:"unitPrice"`
	Amount    string `json:"amount"`
	Operator  string `json:"operator"`
	Location  string `json:"location"`
	Remark    string `json:"remark"`
}

type ContractCreateReq struct {
	ProjectName  string            `json:"projectName" binding:"required"`
	OrderNo      string            `json:"orderNo"`
	OrderDate    string            `json:"orderDate"`
	FromCompany  string            `json:"fromCompany"`
	ToCompany    string            `json:"toCompany"`
	Buyer        string            `json:"buyer"`
	Attn         string            `json:"attn"`
	BuyerEmail   string            `json:"buyerEmail"`
	BuyerTel     string            `json:"buyerTel"`
	AttnTel      string            `json:"attnTel"`
	TotalAmount  string            `json:"totalAmount"`
	DeliveryAddr string            `json:"deliveryAddr"`
	Remark       string            `json:"remark"`
	Items        []ContractItemReq `json:"items"`
}

type ContractUpdateReq struct {
	ProjectName  string            `json:"projectName"`
	OrderNo      string            `json:"orderNo"`
	OrderDate    string            `json:"orderDate"`
	FromCompany  string            `json:"fromCompany"`
	ToCompany    string            `json:"toCompany"`
	Buyer        string            `json:"buyer"`
	Attn         string            `json:"attn"`
	BuyerEmail   string            `json:"buyerEmail"`
	BuyerTel     string            `json:"buyerTel"`
	AttnTel      string            `json:"attnTel"`
	TotalAmount  string            `json:"totalAmount"`
	DeliveryAddr string            `json:"deliveryAddr"`
	Remark       string            `json:"remark"`
	Items        []ContractItemReq `json:"items"`
}

type ContractListReq struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
	OrderNo  string `form:"orderNo"`
}

type ContractDetail struct {
	Contract
	Items []ContractItem `json:"items"`
}
