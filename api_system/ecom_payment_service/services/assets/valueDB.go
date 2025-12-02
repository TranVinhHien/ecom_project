package assets_services

// accoutsTable
const (
	AccoutsTable_ActiveStatus_Active   = "Active"
	AccoutsTable_ActiveStatus_Inactive = "Inactive"
)

// customerTable
const (
	CustomersTable_Gender_Nam = "Nam"
	CustomersTable_Gender_Nu  = "Nữ"
)

// customerTable
const (
	ProductSPUTable_StockStatus_InStock    = "InStock"
	ProductSPUTable_StockStatus_OutOfStock = "OutOfStock"
)
const (
	ProductSPUTable_DeleteStatus_Active  = "Active"
	ProductSPUTable_DeleteStatus_Deleted = "Deleted"
)

// orderTable
const (
	OrderTable_PaymentStatus_ChoThanhToan      = "Chờ Thanh Toán"
	OrderTable_PaymentStatus_ThanhToanTrucTiep = "Thanh Toán Trực Tiếp"
	OrderTable_PaymentStatus_DaThanhToan       = "Đã Thanh Toán"
	OrderTable_PaymentStatus_ThanhToanHetHang  = "Thanh Toán Hết Hạn"
)
const (
	OrderTable_OrderStatus_DaHuy      = "Đã Hủy"
	OrderTable_OrderStatus_ChoXacNhan = "Chờ Xác Nhận"
	OrderTable_OrderStatus_DaXacNhan  = "Đã Xác Nhận"
	OrderTable_OrderStatus_DaGiaoHang = "Đã Giao Hàng"
)
