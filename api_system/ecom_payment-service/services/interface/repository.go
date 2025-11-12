package iservices

// type UserRepository interface {
// 	GetCustomer(ctx context.Context, customer_id string) (services.Customers, error)
// 	GetCustomerByAccountID(ctx context.Context, accoutn_id string) (services.Customers, error)

// 	UpdateCustomers(ctx context.Context, user services.Customers, fn func() error) error

// 	GetAccountByUserName(ctx context.Context, username string) (services.Accounts, error)
// 	UpdateAccount(ctx context.Context, account services.Accounts) error

// 	Register(ctx context.Context, account *services.Accounts, userInfo *services.Customers, roleID string) error
// 	Login(ctx context.Context, username string) (account services.Accounts, role string, err error)

// 	ListCustomerAddresses(ctx context.Context, customer_id string) (addresss []services.CustomerAddress, err error)
// 	CustomerAddresses(ctx context.Context, customer_id, address_id string) (info services.Customers, address services.CustomerAddress, err error)
// 	CreateCustomerAddresses(ctx context.Context, customer_id string, address *services.CustomerAddress) (err error)
// 	UpdateCustomerAddresses(ctx context.Context, customer_id string, address *services.CustomerAddress) (err error)
// 	DeleteCustomerAddresses(ctx context.Context, customer_id string, address_id string) (err error)

// 	UpdateDeviceRegistrationToken(ctx context.Context, customer_id string, token string) (err error)
// }
// type CategoriesRepository interface {
// 	ListCategories(ctx context.Context) ([]services.Categorys, error)
// 	ListCategoriesByID(ctx context.Context, cate_id string) ([]services.Categorys, error)
// }

// type ĐiscountRepository interface {
// 	ListDiscount(ctx context.Context, query services.QueryFilter) (is []services.Discounts, totalPages, totalElements int, err error)
// 	Discount(ctx context.Context, discount string) (i services.Discounts, err error)
// 	GetDiscountForNoti(ctx context.Context) (is []services.Discounts, err error)
// }
// type PaymentRepository interface {
// 	ListPayment(ctx context.Context) (is []services.PaymentMethods, err error)
// }
// type OrderRepository interface {
// 	TXCreateOrdder(ctx context.Context, order *services.Orders, orderDetail []services.OrderDetail) (err error)
// 	UpdateOrder(ctx context.Context, order services.Orders) (err error)
// 	GetOrdersByUserID(ctx context.Context, userID string, query services.QueryFilter) (items []services.Orders, totalPages, totalElements int, err error)
// 	GetOrderDetailByOrderIDs(ctx context.Context, orderIDs []string) (is []services.OrderDetail, err error)
// 	GetOrderByID(ctx context.Context, orderID string) (i services.Orders, err error)
// 	TXCancelOrder(ctx context.Context, orderID string) error
// 	CheckUserOrder(ctx context.Context, userID, products_spu_id string) (count int64, err error)
// }

// type ProductRepository interface {
// 	GetProductsBySKUs(ctx context.Context, product_sku_ids []string) (is []services.ProductSkusDetail, err error)
// }

// type RatingRepository interface {
// 	GetRatings(ctx context.Context, query services.QueryFilter) (items []services.Ratings, totalPages, totalElements int, err error)
// 	CreateRating(ctx context.Context, rating services.Ratings) (err error)
// }
// type ProductsRepository interface {
// 	GetAllProductSimple(ctx context.Context, query services.QueryFilter) (items []services.ProductSimple, totalPages, totalElements int, err error)
// 	GetProductDetail(ctx context.Context, productSpuID string) (product_detail services.ProductDetail, err error)
// }
