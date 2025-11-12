package services_assets_sendMessage

import (
	"fmt"

	"firebase.google.com/go/messaging"
)

func ThanhToanThanhCong(orderID string, total float64) *messaging.Notification {
	return &messaging.Notification{
		Title: "Thanh toán thành công",
		Body:  "Bạn đã thanh toán thành công đơn hàng." + orderID + "với giá trị" + fmt.Sprint("%f", total) + "VNĐ",
		// ImageURL: "https://f6e9-118-68-56-216.ngrok-free.app/v1/media/products?id=images/phu-kien-thoi-trang___phu-kien-nu___phu-kien-nu-khac/100477499_1.jpg",
	}
}

// func MaGiamGiaMoi(discount []services.Discounts) *messaging.Notification {
// 	max := discount[0].DiscountValue
// 	for _, dis := range discount {
// 		if dis.DiscountValue > max {
// 			max = dis.DiscountValue
// 		}
// 	}
// 	return &messaging.Notification{
// 		Title: "Mã giảm giá",
// 		Body:  "Có " + fmt.Sprint(" ", len(discount)) + " mã giảm giá sẽ được sử dụng trong 1 tiếng tới với giá trị giảm giá lớn nhất lên đến " + fmt.Sprint(" ", max) + "vnđ",
// 		// ImageURL: "https://f6e9-118-68-56-216.ngrok-free.app/v1/media/products?id=images/phu-kien-thoi-trang___phu-kien-nu___phu-kien-nu-khac/100477499_1.jpg",
// 	}
// }
