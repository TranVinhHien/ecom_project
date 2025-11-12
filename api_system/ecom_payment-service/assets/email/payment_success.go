package email

import (
	"bytes"
	"fmt"
	"html/template"

	services "github.com/TranVinhHien/ecom_payment-service/services/entity"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ====================================================================
// CÁC STRUCT ĐẦU VÀO (Theo định nghĩa của bạn)
// ====================================================================

// ====================================================================
// HÀM HELPER (Định dạng tiền tệ)
// ====================================================================

// formatCurrency là một helper để hiển thị tiền tệ theo chuẩn Việt Nam.
func formatCurrency(amount interface{}) string {
	var f float64 // Biến tạm để lưu giá trị đã ép kiểu

	// Sử dụng switch để kiểm tra kiểu dữ liệu
	switch v := amount.(type) {
	case int:
		f = float64(v)
	case int64:
		f = float64(v)
	case float32:
		f = float64(v)
	case float64:
		f = v
	default:
		// Nếu không phải kiểu số, trả về chuỗi rỗng hoặc báo lỗi
		return ""
	}

	// Sử dụng message printer cho ngôn ngữ tiếng Việt
	p := message.NewPrinter(language.Vietnamese)
	// Định dạng số với 0 chữ số thập phân, vì là VND
	return p.Sprintf("%.0f ₫", f)
}

// ====================================================================
// HÀM CHÍNH (Render Template)
// ====================================================================

// GeneratePaymentSuccessEmail là hàm nghiệp vụ trong Notification Service.
// Nó nhận payload, render template HTML, và trả về nội dung email.
func GeneratePaymentSuccessEmail(data services.CombinedDataPayLoadMoMo) (string, error) {

	// Định nghĩa template HTML (giống hệt mẫu ở trên)
	// Trong production, bạn nên đọc template này từ file .html
	const templateStr = `
<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Xác nhận Thanh toán Đơn hàng</title>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@400;600;700&display=swap');
        body{margin:0;padding:0;background-color:#f4f1e9;font-family:'Cormorant Garamond', Georgia, 'Times New Roman', serif;color:#3a3a3a;}
        .container{width:90%;max-width:680px;margin:20px auto;background-color:#fdfaf2;border:1px solid #dcd1b3;border-radius:8px;box-shadow:0 5px 15px rgba(0, 0, 0, 0.05);border:1px solid #c8b77b;position:relative;padding:5px;}
        .inner-border{border:2px double #d4af37;padding:25px;border-radius:4px;}
        .header{text-align:center;padding-bottom:20px;border-bottom:2px solid #e8e1c5;}
        .header h1{font-size:36px;font-weight:600;color:#2c3e50;margin:0;}
        .seal{display:block;width:60px;height:60px;background-color:#8c0000;border-radius:50%;margin:15px auto 10px;color:#d4af37;font-size:32px;font-weight:bold;line-height:60px;text-align:center;font-family:'Times New Roman', serif;box-shadow:0 0 5px rgba(0,0,0,0.5), inset 0 1px 3px rgba(255,255,255,0.3);position:relative;}
        .seal::before{content:'';position:absolute;top:2px;left:2px;right:2px;bottom:2px;border:1px dashed #d4af37;border-radius:50%;}
        .content{padding:25px 0;line-height:1.7;font-size:18px;}
        .content p{margin-bottom:20px;}
        .order-details{margin-top:20px;border-top:1px solid #e8e1c5;padding-top:20px;}
        .order-details h2{font-size:24px;font-weight:600;color:#2c3e50;border-bottom:1px solid #e8e1c5;padding-bottom:10px;margin-bottom:20px;}
        .info-table{width:100%;margin-bottom:20px;}
        .info-table td{padding:5px 0;font-size:17px;}
        .info-table td:first-child{font-weight:600;color:#574f39;width:180px;}
        .item-list{width:100%;border-collapse:collapse;margin-top:15px;}
        .item-list th, .item-list td{padding:12px 10px;text-align:left;border-bottom:1px solid #e8e1c5;}
        .item-list th{background-color:#fdfaf2;font-size:16px;font-weight:700;color:#574f39;}
        .item-list .item-image{width:60px;height:60px;object-fit:cover;border-radius:4px;border:1px solid #e8e1c5;}
        .item-list .item-name{font-size:17px;font-weight:600;}
        .item-list .item-qty{font-size:16px;text-align:center;}
        .item-list .item-price{font-size:17px;font-weight:600;text-align:right;white-space:nowrap;}
        .total-section{margin-top:25px;text-align:right;}
        .total-section p{font-size:22px;font-weight:700;color:#333;margin:5px 0;}
        .total-section .grand-total{font-size:28px;color:#8c0000;}
        .footer{margin-top:30px;padding-top:20px;border-top:2px solid #e8e1c5;text-align:center;font-size:14px;color:#777;font-style:italic;}
    </style>
</head>
<body>
    <div class="container">
        <div class="inner-border">
            <div class="header">
                <div class="seal">Y</div>
                <h1>Thanh Toán Thành Công</h1>
            </div>

            <div class="content">
                <p>Kính thưa Quý khách {{.Info.Name}},</p>
                <p>Chúng tôi trân trọng thông báo rằng Giao dịch thanh toán cho đơn hàng của Quý khách đã được xử lý thành công. Chúng tôi xin bày tỏ lòng biết ơn sâu sắc vì Quý khách đã tin tưởng và lựa chọn dịch vụ của chúng tôi.</p>
                <p>Đơn hàng của Quý khách hiện đang được chuẩn bị và sẽ sớm được vận chuyển. Quý khách có thể theo dõi trạng thái đơn hàng trong tài khoản của mình.</p>

                <div class="order-details">
                    <h2>Tóm tắt Đơn hàng</h2>
                    
                    <table class="info-table">
                        <tr>
                            <td>Mã Đơn hàng:</td>
                            <td>{{.Order.OrderID}}</td>
                        </tr>
                        <tr>
                            <td>Mã Giao dịch:</td>
                            <td>{{.TransactionID}}</td>
                        </tr>
                        <tr>
                            <td>Giao đến:</td>
                            <td>{{.Info.Name}}</td>
                        </tr>
                        <tr>
                            <td>Địa chỉ:</td>
                            <td>{{.Info.Address}}</td>
                        </tr>
                         <tr>
                            <td>Điện thoại:</td>
                            <td>{{.Info.PhoneNumber}}</td>
                        </tr>
                    </table>
                    
                    <h2 style="margin-top: 30px;">Chi tiết Sản phẩm</h2>
                    <table class="item-list">
                        <thead>
                            <tr>
                                <th colspan="2">Sản phẩm</th>
                                <th style="text-align: center;">Số lượng</th>
                                <th style="text-align: right;">Tổng cộng</th>
                            </tr>
                        </thead>
                        <tbody>
                            {{range .Items}}
                            <tr>
                                <td style="width: 60px;">
                                    <img src="{{.ImageURL}}" alt="{{.Name}}" class="item-image" style="width: 60px; height: 60px; object-fit: cover; border-radius: 4px; border: 1px solid #e8e1c5;">
                                </td>
                                <td>
                                    <div class="item-name">{{.Name}}</div>
                                </td>
                                <td>
                                    <div class="item-qty">{{.Quantity}}</div>
                                </td>
                                <td>
                                    <div class="item-price">{{formatCurrency .TotalPrice}}</div>
                                </td>
                            </tr>
                            {{end}}
                        </tbody>
                    </table>

                    <div class="total-section">
                        <p>Tổng cộng:</p>
                        <p class="grand-total">{{formatCurrency .Order.TotalAmount}}</p>
                    </div>
                </div>
            </div>

            <div class="footer">
                <p>Trân trọng cảm ơn Quý khách.</p>
                <p>&copy; 2025 Sàn Thương mại điện tử LemarCheNoble. Bảo lưu mọi quyền.</p>
            </div>
        </div>
    </div>
</body>
</html>
`

	// Đăng ký hàm helper vào template
	// Sử dụng text/template để parse hàm, sau đó dùng html/template để thực thi an toàn
	// Cập nhật: html/template.New().Funcs() là cách làm đúng
	tmpl, err := template.New("paymentSuccessEmail").
		Funcs(template.FuncMap{
			"formatCurrency": formatCurrency,
		}).
		Parse(templateStr)

	if err != nil {
		// Lỗi này là lỗi hệ thống (lỗi parse template), cần log chi tiết
		return "", fmt.Errorf("lỗi parse template email: %w", err)
	}

	// Tạo một buffer để ghi kết quả HTML
	var renderedEmail bytes.Buffer

	// Thực thi template với dữ liệu
	if err := tmpl.Execute(&renderedEmail, data); err != nil {
		// Lỗi này có thể do dữ liệu (data)
		return "", fmt.Errorf("lỗi render template email: %w", err)
	}

	// Trả về chuỗi HTML đã render
	// Bước tiếp theo (trong service) sẽ là:
	// 1. Inline CSS
	// 2. Lấy email người dùng từ User Service (dựa trên user_id)
	// 3. Gọi Email Provider API (như trong Postman )
	return renderedEmail.String(), nil
}
