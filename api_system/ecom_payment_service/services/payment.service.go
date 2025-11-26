package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	config_assets "github.com/TranVinhHien/ecom_payment_service/assets/config"
	"github.com/TranVinhHien/ecom_payment_service/assets/email"
	db "github.com/TranVinhHien/ecom_payment_service/db/sqlc"
	assets_services "github.com/TranVinhHien/ecom_payment_service/services/assets"
	entity "github.com/TranVinhHien/ecom_payment_service/services/entity"

	// services "github.com/TranVinhHien/ecom_payment_service/services/entity"
	"github.com/google/uuid"
)

func (s *service) PaymentMethodDetail(ctx context.Context, id string) (map[string]interface{}, *assets_services.ServiceError) {

	payment, err := s.repository.GetPaymentMethodByID(ctx, id)
	if err != nil {
		// fmt.Println("Error GetPaymentMethodByID:", err)
		return nil, assets_services.NewError(400, err)
	}
	result, err := assets_services.HideFields(payment, "data")
	if err != nil {
		// fmt.Println("Error HideFields:", err)
		return nil, assets_services.NewError(400, err)
	}
	return result, nil
}
func (s *service) ListPaymentMethod(ctx context.Context) (map[string]interface{}, *assets_services.ServiceError) {

	payments, err := s.repository.ListActivePaymentMethods(ctx)
	if err != nil {
		// fmt.Println("Error ListPayment:", err)
		return nil, assets_services.NewError(400, err)
	}
	result, err := assets_services.HideFields(payments, "data")
	if err != nil {
		// fmt.Println("Error HideFields:", err)
		return nil, assets_services.NewError(400, err)
	}
	return result, nil
}
func (s *service) GetPaymentMethodByID(ctx context.Context, id string) (map[string]interface{}, *assets_services.ServiceError) {

	payment, err := s.repository.GetPaymentMethodByID(ctx, id)
	if err != nil {
		fmt.Println("Error GetPaymentMethodByID:", err)
		return nil, assets_services.NewError(400, err)
	}
	result, err := assets_services.HideFields(payment, "data")
	if err != nil {
		fmt.Println("Error HideFields:", err)
		return nil, assets_services.NewError(400, err)
	}
	return result, nil
}

func (s *service) InitPayment(ctx context.Context, userId string, email string, order entity.InitPaymentParams) (map[string]interface{}, *assets_services.ServiceError) {
	payment, err := s.repository.GetPaymentMethodByID(ctx, order.PaymentMethodID) // Sửa tên trường thành PaymentMethodID
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, assets_services.NewError(404, fmt.Errorf("phương thức thanh toán không hợp lệ"))
		}
		return nil, assets_services.NewError(500, err) // Lỗi DB
	}
	// order.UserInfo.Address =// email           // gán email vào address
	var paymentResult map[string]interface{} // Lưu kết quả trả về (URL hoặc thông báo COD)
	transactionID := ""                      // Lưu transactionID để dùng sau

	err = s.repository.ExecTS(ctx, func(tx db.Querier) error {
		// --- Bước 1: Tạo Transaction (Luôn là PENDING ban đầu) ---
		transactionID = uuid.New().String()
		err := tx.CreateTransaction(ctx, db.CreateTransactionParams{
			ID:              transactionID,
			TransactionCode: generateOrderCode(),
			OrderID:         sql.NullString{String: order.OrderID, Valid: true},
			PaymentMethodID: payment.ID,
			Amount:          fmt.Sprintf("%.2f", order.Amount), // Cẩn thận với float/string
			Currency:        "VND",
			Type:            db.TransactionsTypePAYMENT,
			Status:          db.TransactionsStatusPENDING,             // Luôn PENDING ban đầu
			Notes:           sql.NullString{String: "", Valid: false}, // Hoặc có thể thêm ghi chú nếu cần
		})
		if err != nil {
			return fmt.Errorf("lỗi tạo transaction: %w", err)
		}

		// // ✅ GỬI KAFKA EVENT: Transaction Created
		// go func() {
		// 	eventData := map[string]interface{}{
		// 		"order_id":         order.OrderID,
		// 		"payment_method":   payment.Code,
		// 		"amount":           order.Amount,
		// 		"currency":         "VND",
		// 	}
		// 	if err := s.eventPublisher.PublishTransactionCreated(context.Background(), transactionID, userId, eventData); err != nil {
		// 		fmt.Println("❌ Lỗi gửi Kafka event TransactionCreated:", err)
		// 	}
		// }()

		// --- Bước 2: Tạo bản ghi chi phí Sàn ---
		// Lưu ý: Đảm bảo bạn đã có hàm CreateOrderPlatformCost trong sqlc
		err = tx.CreateOrderPlatformCost(ctx, db.CreateOrderPlatformCostParams{
			OrderID:                        order.OrderID,
			PaymentTransactionID:           transactionID,
			SiteOrderVoucherDiscountAmount: fmt.Sprintf("%.2f", order.SiteOrderVoucherDiscountAmount),
			SitePromotionDiscountAmount:    fmt.Sprintf("%.2f", order.SitePromotionDiscountAmount),
			SiteShippingDiscountAmount:     fmt.Sprintf("%.2f", order.SiteShippingDiscountAmount),
			TotalSiteFundedProductDiscount: fmt.Sprintf("%.2f", order.TotalSiteFundedProductDiscount),
		})
		if err != nil {
			return fmt.Errorf("lỗi tạo platform cost: %w", err)
		}

		// --- Bước 3: Tạo các bản ghi Shop Order Settlements ---
		// Lưu ý: Đảm bảo hàm sqlc CreateShopOrderSettlement đã được cập nhật
		// để KHÔNG còn các trường site_shipping_discount, site_order_discount
		for _, item := range order.SettlementDetails {
			err = tx.CreateShopOrderSettlement(ctx, db.CreateShopOrderSettlementParams{
				ID:                        uuid.New().String(),
				ShopOrderID:               item.ShopOrderID,
				OrderTransactionID:        transactionID,
				Status:                    db.ShopOrderSettlementsStatusPENDINGSETTLEMENT,
				OrderSubtotal:             fmt.Sprintf("%.2f", item.OrderSubtotal),
				ShopFundedProductDiscount: fmt.Sprintf("%.2f", item.ShopFundedProductDiscount),
				SiteFundedProductDiscount: fmt.Sprintf("%.2f", item.SiteFundedProductDiscount),
				ShopVoucherDiscount:       fmt.Sprintf("%.2f", item.ShopVoucherDiscount),
				SiteOrderDiscount:         fmt.Sprintf("%.2f", item.SiteOrderDiscount),
				SiteShippingDiscount:      fmt.Sprintf("%.2f", item.SiteShippingDiscount),
				ShippingFee:               fmt.Sprintf("%.2f", item.ShippingFee),
				ShopShippingDiscount:      fmt.Sprintf("%.2f", item.ShopShippingDiscount),
				CommissionFee:             fmt.Sprintf("%.2f", item.CommissionFee),
				NetSettledAmount:          fmt.Sprintf("%.2f", item.NetSettledAmount),
			})
			if err != nil {
				return fmt.Errorf("lỗi tạo settlement cho %s: %w", item.ShopOrderID, err)
			}
		}

		// --- Bước 4: Khởi tạo luồng thanh toán (Nếu Online) ---
		// KHÔNG thực hiện hạch toán kế toán ở đây
		if payment.Type == db.PaymentMethodsTypeONLINE {
			// thực hiện tạo url thanh toán trong transaction
			if payment.Code == "MOMO" {
				orderDetailMOMO := make([]entity.Product_MOMO, 0)
				for _, item := range order.Items {
					orderDetailMOMO = append(orderDetailMOMO, entity.Product_MOMO{
						ID:         item.ProductID,
						Name:       item.Name,
						ImageURL:   item.ImageURL,
						Price:      int64(item.Price),
						Currency:   "VND",
						Quantity:   item.Quantity,
						TotalPrice: int64(item.Price) * int64(item.Quantity),
					})
				}

				payloadParamsMoMo := entity.CombinedDataPayLoadMoMo{
					Info:          order.UserInfo,
					Order:         entity.OrderValue{OrderID: order.OrderID, TotalAmount: order.Amount},
					Items:         orderDetailMOMO,
					TransactionID: transactionID,
					Email:         email,
				}
				payload := createMoMoPayload(s.env, payloadParamsMoMo)
				defer s.redis.AddTransactionOnline(ctx, userId, payloadParamsMoMo, s.env.OrderDuration+10*time.Minute) // Lưu thông tin đơn hàng trong Redis để dùng khi callback

				// log đầy đủ thông tin của payload dễ nhìn
				urlCallBack, errors := callMoMoGetURL(s.env, payload)
				if errors != nil {
					return errors.Err
				}
				paymentResult = urlCallBack
			} else {
				return fmt.Errorf("phương thức thanh toán online không được hỗ trợ: %s", payment.Code)
			}

		} else { // Xử lý Offline (COD)
			// Không cần gọi cổng thanh toán
			// Có thể chuẩn bị một thông báo thành công cho COD
			paymentResult = map[string]interface{}{
				"message":    "Đơn hàng COD đã được tạo thành công.",
				"order_code": order.OrderID, // Hoặc mã order code thân thiện hơn
				"amount":     order.Amount,
			}
		}

		return nil // Commit transaction DB nếu không có lỗi
	})

	if err != nil {
		// Lỗi xảy ra trong transaction (DB hoặc MoMo API...)
		return nil, assets_services.NewError(400, err)
	}
	result, _ := assets_services.HideFields(paymentResult, "data")

	return result, nil
}

func (s *service) CallBackMoMo(ctx context.Context, tran entity.TransactionMoMO) {
	// su ly api thanh cong
	// sử lý thoong báo tại đây
	if tran.ResultCode == 0 {
		fmt.Println("thanh toan momo thanh cong,", tran)
		err := s.repository.ExecTS(ctx, func(tx db.Querier) error {

			// get transaction
			transactionDB, err := tx.GetTransactionByID(ctx, tran.RequestID)
			if err != nil {

				return fmt.Errorf("lỗi khi lấy transactionid: %s", err.Error())
			}
			if transactionDB.Amount != fmt.Sprintf("%0.2f", float64(tran.Amount)) {

				return fmt.Errorf("số tiền không khớp: %s", transactionDB.Amount)
			}
			// cập nhật trạng thái transaction
			err = tx.UpdateTransactionStatus(ctx, db.UpdateTransactionStatusParams{
				ID:                   transactionDB.ID,
				Status:               db.TransactionsStatusSUCCESS,
				ProcessedAt:          sql.NullTime{Time: time.Now(), Valid: true},
				GatewayTransactionID: sql.NullString{String: fmt.Sprint(tran.TransactionID), Valid: true},
				Notes:                sql.NullString{String: tran.Message, Valid: true},
			})
			if err != nil {
				return err
			}
			// tạo ledger_entries
			err = tx.CreateLedgerEntry(ctx, db.CreateLedgerEntryParams{
				LedgerID:      s.env.PlatformID,
				TransactionID: tran.RequestID,
				Amount:        fmt.Sprintf("%.2f", float64(tran.Amount)),
				Type:          db.LedgerEntriesTypeCREDIT,
				Description:   fmt.Sprintf("Nạp tiền từ MoMo, transactionID: %s", tran.RequestID),
			})
			if err != nil {
				return err
			}
			// cập nhật số tiền hold trong account_ledgers
			err = tx.UpdateLedgerBalances(ctx, db.UpdateLedgerBalancesParams{
				ID:                   s.env.PlatformID,
				BalanceChange:        fmt.Sprintf("%.2f", float64(0)),
				PendingBalanceChange: fmt.Sprintf("%.2f", float64(tran.Amount)),
			})
			if err != nil {
				return err
			}
			//SỬ LÝ Ở CUỐI
			orderID := tran.RequestID
			// Gọi Kafka sử lý khi đơn hàng thanh toán thành công:
			// err = s.apiServer.UpdateOrderCallback(tran.OrderID)
			// if err != nil {
			// 	return err
			// }

			// ✅ GỬI KAFKA EVENT: Payment Completed
			// gửi tới serivce order để cập nhật trạng thái đơn hàng
			// gửi sự kiện tới email service để gửi email thanh toán thành công cho khách hàng

			// lấy thông tin transaction để gửi event
			orderInfo, err := s.redis.GetTransactionOnlineWithIDTran(ctx, orderID)
			if err != nil {
				fmt.Println("❌ Lỗi lấy thông tin transaction:", err)
				return err
			}
			if orderInfo == nil {
				fmt.Println("❌ Không tìm thấy thông tin đơn hàng cho transaction:", orderID)
				return fmt.Errorf("không tìm thấy thông tin đơn hàng cho transaction: %s", orderID)
			}
			// init content email
			emailContent, err := email.GeneratePaymentSuccessEmail(*orderInfo)
			if err != nil {
				fmt.Println("❌ Lỗi tạo nội dung email:", err)
				return err
			}

			go func() {

				// create event data
				//add data send mail
				eventData := map[string]interface{}{
					"order_id":       transactionDB.OrderID.String,
					"transaction_id": transactionDB.ID,
					"amount":         tran.Amount,
					"payment_method": "MOMO",
					"message":        tran.Message,
				}
				if err := s.producer.PaymentCompleted(context.Background(), transactionDB.OrderID.String, eventData); err != nil {
					fmt.Println("❌ Lỗi gửi Kafka event PaymentCompleted:", err)
				}
			}()

			err = s.email.SendEmail(ctx,
				orderInfo.Email,
				orderInfo.Info.Name,
				"Xác nhận thanh toán thành công đơn hàng "+orderInfo.Order.OrderID,
				emailContent,
			)
			if err != nil {
				fmt.Println("Lỗi gửi email thanh toán thành công:", err)
				return err
			}
			s.redis.DeleteTransactionOnline(ctx, orderID)
			return nil
		})
		// sử lý nếu gặp lỗi phải lưu transaction vào 1 nơi để thực hiện sử lý lại sau 1 khoảng thời gian.
		if err != nil {
			fmt.Println("Lỗi cập nhật transaction sau khi thanh toán MoMo:", err)
		}
	} else {
		fmt.Println("thanh toan momo that bai,", tran)
	}
}
func (s *service) GetURLOrderMoMOAgain(ctx context.Context, user_id string) (map[string]interface{}, *assets_services.ServiceError) {
	// check neu co order thi khong cho nguoi dung thanh toan
	payloadParamsMoMo, err := s.redis.GetTransactionOnline(ctx, user_id)
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("error redis.GetOrderOnline:  %s ", err.Error()))
	}
	payload := createMoMoPayload(s.env, *payloadParamsMoMo)
	return callMoMoGetURL(s.env, payload)
}

func createMoMoPayload(env config_assets.ReadENV, payloadd entity.CombinedDataPayLoadMoMo) entity.Payload_MOMO {

	var partnerCode = "MOMO"            // mã code của MOMO
	var extraData = ""                  // ko biết
	var partnerName = "Le Marché Noble" // tên đối tác
	var storeId = "MoMoTestStore"       // để đại chứ không hiểu :Mã cửa hàng
	var orderGroupId = ""               // orderGroupId được MoMo cung cấp để phân nhóm đơn hàng cho các hoạt động vận hành sau này. Vui lòng liên hệ với MoMo để biết chi tiết cách sử dụng
	var autoCapture = true              //  Nếu giá trị false, giao dịch sẽ không tự động capture. Mặc định là true
	var lang = "vi"
	var requestType = "payWithMethod"                                                   // captureWallet
	var amount = strconv.Itoa(int(payloadd.Order.TotalAmount))                          // Số tiền thanh toán Số tiền cần thanh toán// Nhỏ Nhất: 1.000 VND// Tối đa: 50.000.000 VND// Tiền tệ: VND// Kiểu dữ liệu: Long
	var orderId = payloadd.Order.OrderID                                                // mã đơn hàng
	var orderInfo = fmt.Sprintf("Thanh toán %s VNĐ cho đơn hàng : %s", amount, orderId) // thông tin nhắn gửi
	var accessKey = env.AccessKeyMoMo                                                   //
	var secretKey = env.SecretKeyMoMo
	var redirectUrl = env.RedirectURL      // Một URL của đối tác. URL này được sử dụng để chuyển trang (redirect) từ MoMo về trang mua hàng của đối tác sau khi khách hàng thanh toán.Hỗ trợ: AppLink and WebLink
	var ipnUrl = env.IpnURL                // API của đối tác. Được MoMo sử dụng để gửi kết quả thanh toán theo phương thức IPN (server-to-server)
	var requestId = payloadd.TransactionID // Định danh duy nhất cho mỗi yêu cầu Đối tác sử dụng requestId để xử lý

	var rawSignature bytes.Buffer
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(accessKey)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(amount)
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(extraData)
	rawSignature.WriteString("&ipnUrl=")
	rawSignature.WriteString(ipnUrl)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(orderId)
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(orderInfo)
	rawSignature.WriteString("&partnerCode=")
	rawSignature.WriteString(partnerCode)
	// chuyển hướng về ứng dụng của mình
	rawSignature.WriteString("&redirectUrl=")
	rawSignature.WriteString(redirectUrl)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(requestId)
	rawSignature.WriteString("&requestType=")
	rawSignature.WriteString(requestType)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	hmac := hmac.New(sha256.New, []byte(secretKey))

	// Write Data to it
	hmac.Write(rawSignature.Bytes())

	// Get result and encode as hexadecimal string
	signature := hex.EncodeToString(hmac.Sum(nil))

	var payload = entity.Payload_MOMO{
		PartnerCode:  partnerCode,
		AccessKey:    accessKey,
		RequestID:    requestId,
		Amount:       amount,
		RequestType:  requestType,
		RedirectUrl:  redirectUrl,
		IpnUrl:       ipnUrl,
		OrderID:      orderId,
		StoreId:      storeId,
		PartnerName:  partnerName,
		OrderGroupId: orderGroupId,
		AutoCapture:  autoCapture,
		Lang:         lang,
		OrderInfo:    orderInfo,
		ExtraData:    extraData,
		Signature:    signature,
		Items:        payloadd.Items,
		UserInfo: entity.User_MOMO{
			Name:        payloadd.Info.Name,
			PhoneNumber: payloadd.Info.PhoneNumber,
			Address:     payloadd.Info.Address,
		},
	}

	return payload
}

func callMoMoGetURL(env config_assets.ReadENV, payload entity.Payload_MOMO) (map[string]interface{}, *assets_services.ServiceError) {

	var jsonPayload []byte
	var err error
	jsonPayload, err = json.Marshal(payload)
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("error when json.Marshal %s", err.Error()))
	}
	//send HTTP to momo endpoint
	resp, err := http.Post(env.EndPointMoMo, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("error when send HTTP to momo endpoint: %s", err.Error()))
	}

	//result
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
func generateOrderCode() string {
	timestamp := time.Now().Format("20060102")
	randomPart := uuid.New().String()[:8]
	return fmt.Sprintf("YAN%s%s", timestamp, randomPart)
}

// func groupItemsByShop(items []entity.DetailItem) map[string][]entity.DetailItem {
// 	shopMap := make(map[string][]entity.DetailItem)
// 	for _, item := range items {
// 		shopMap[item.ShopID] = append(shopMap[item.ShopID], item)
// 	}
// 	return shopMap
// }
