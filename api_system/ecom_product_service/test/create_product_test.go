package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	baseURL = "http://172.26.127.95:9001/v1/product/create"
	// token   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJoaWVubGF6YWRhIiwic2NvcGUiOiJTWVNURU0iLCJpc3MiOiJsZW1hcmNoZW5vYmxlLmlkLnZuIiwiZXhwIjo0OTE3NTExMjUyLCJpYXQiOjE3NjE3NTEyNTIsInVzZXJJZCI6IjE2NzQwOGUzLWFmZWYtNDhiOS04ZTRmLTZkZDQxZWJmMzQ2NCIsImp0aSI6ImU2YzgyN2E2LTIyOTYtNGNlOC1iMjQ1LWM3MDIxNWM4MGJjNyIsImVtYWlsIjoidmluaGhpZW4xMnpAZ21haWwuY29tIn0.CPnP_NqB_WtaQb9X43YKFav8wYzdqB14jFNtnPr74as"
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJoaWVubGF6YWRhIiwic2NvcGUiOiJTRUxMRVIiLCJpc3MiOiJsZW1hcmNoZW5vYmxlLmlkLnZuIiwiZXhwIjo0OTE3NTExMjUyLCJpYXQiOjE3NjE3NTEyNTIsInVzZXJJZCI6IjE2NzQwOGUzLWFmZWYtNDhiOS04ZTRmLTZkZDQxZWJmMzQ2NCIsImp0aSI6ImU2YzgyN2E2LTIyOTYtNGNlOC1iMjQ1LWM3MDIxNWM4MGJjNyIsImVtYWlsIjoidmluaGhpZW4xMnpAZ21haWwuY29tIn0.2N5eJfouPOkZ47Nh2PLdrtJ_Md7zoDuvCfPn7XqncuQ"
)

// ProductRequest represents the product creation request structure
type ProductRequest struct {
	Name                      string        `json:"name"`
	Key                       string        `json:"key"`
	Description               string        `json:"description"`
	ShortDescription          string        `json:"short_description"`
	BrandID                   string        `json:"brand_id"`
	CategoryID                string        `json:"category_id"`
	ShopID                    string        `json:"shop_id"`
	ProductIsPermissionReturn bool          `json:"product_is_permission_return"`
	ProductIsPermissionCheck  bool          `json:"product_is_permission_check"`
	OptionValue               []OptionValue `json:"option_value"`
	ProductSKU                []ProductSKU  `json:"product_sku"`
}

type OptionValue struct {
	OptionName string `json:"option_name"`
	Value      string `json:"value"`
}

type ProductSKU struct {
	SkuCode     string        `json:"sku_code"`
	Price       float64       `json:"price"`
	Quantity    int32         `json:"quantity"`
	Weight      float64       `json:"weight"`
	OptionValue []OptionValue `json:"option_value"`
}

type OptionImageInfo struct {
	OptionName string
	Value      string
	ImagePath  string
}

// TestCreateProduct tests creating a single product
func TestCreateProduct(t *testing.T) {
	// Prepare product data - matching exact structure from request
	productData := ProductRequest{
		Name:                      "Áo thun ngắn tay cúc cài nữ Babyttee 6 cúc in chữ R",
		Key:                       fmt.Sprintf("ao-thun-nam-hien-%s", uuid.New().String()[:8]),
		Description:               "Áo thun nữ kiểu babyttee croptop 6 cúc in chữ R, chất thun tăm co giãn, form ôm tôn dáng. Mềm mịn, thoáng mát, phù hợp đi học, đi chơi, dạo phố.",
		ShortDescription:          "Áo thun tăm croptop nữ 6 cúc chữ R, co giãn thoải mái, tôn dáng.",
		BrandID:                   "3d4877bb-9974-4d1f-87b5-ad801efb99ff",
		CategoryID:                "711e020a-3ffc-4eb6-b0ce-8873106bdf65",
		ShopID:                    "uuid-shop",
		ProductIsPermissionReturn: true,
		ProductIsPermissionCheck:  true,
		OptionValue: []OptionValue{
			{OptionName: "Màu Sắc", Value: "Xanh than chữ R"},
			{OptionName: "Màu Sắc", Value: "Trắng 3 cúc"},
			{OptionName: "Màu Sắc", Value: "Đen 3 cúc"},
			{OptionName: "Màu Sắc", Value: "Trắng 4 cúc"},
			{OptionName: "Màu Sắc", Value: "Legging đùi"},
			{OptionName: "Size", Value: "S (Dưới 50KG)"},
			{OptionName: "Size", Value: "M (Dưới 54KG)"},
		},
		ProductSKU: []ProductSKU{
			{
				SkuCode:  fmt.Sprintf("XANH-THAN-CHU-123R-S-%s", uuid.New().String()[:6]),
				Price:    43000,
				Quantity: 50,
				Weight:   0.25,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: "Xanh than chữ R"},
					{OptionName: "Size", Value: "S (Dưới 50KG)"},
				},
			},
			{
				SkuCode:  fmt.Sprintf("XANH-THAN-CH41U-R-M-%s", uuid.New().String()[:6]),
				Price:    43000,
				Quantity: 50,
				Weight:   0.27,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: "Xanh than chữ R"},
					{OptionName: "Size", Value: "M (Dưới 54KG)"},
				},
			},
			{
				SkuCode:  fmt.Sprintf("TRANG-3-5CUC-S-%s", uuid.New().String()[:6]),
				Price:    43000,
				Quantity: 50,
				Weight:   0.25,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: "Trắng 3 cúc"},
					{OptionName: "Size", Value: "S (Dưới 50KG)"},
				},
			},
			{
				SkuCode:  fmt.Sprintf("TRANG-3-CUgC-M-%s", uuid.New().String()[:6]),
				Price:    43000,
				Quantity: 50,
				Weight:   0.27,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: "Trắng 3 cúc"},
					{OptionName: "Size", Value: "M (Dưới 54KG)"},
				},
			},
			{
				SkuCode:  fmt.Sprintf("DEN-31-CUC-S-%s", uuid.New().String()[:6]),
				Price:    43000,
				Quantity: 50,
				Weight:   0.25,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: "Đen 3 cúc"},
					{OptionName: "Size", Value: "S (Dưới 50KG)"},
				},
			},
			{
				SkuCode:  fmt.Sprintf("DEN-3-3CUC-M-%s", uuid.New().String()[:6]),
				Price:    43000,
				Quantity: 50,
				Weight:   0.27,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: "Đen 3 cúc"},
					{OptionName: "Size", Value: "M (Dưới 54KG)"},
				},
			},
			{
				SkuCode:  fmt.Sprintf("TRAN5G-4-CUC-S-%s", uuid.New().String()[:6]),
				Price:    43000,
				Quantity: 50,
				Weight:   0.25,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: "Trắng 4 cúc"},
					{OptionName: "Size", Value: "S (Dưới 50KG)"},
				},
			},
			{
				SkuCode:  fmt.Sprintf("TRAN1G-4-CUC-M-%s", uuid.New().String()[:6]),
				Price:    43000,
				Quantity: 50,
				Weight:   0.27,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: "Trắng 4 cúc"},
					{OptionName: "Size", Value: "M (Dưới 54KG)"},
				},
			},
			{
				SkuCode:  fmt.Sprintf("LEGGING4-DUI-S-%s", uuid.New().String()[:6]),
				Price:    43000,
				Quantity: 50,
				Weight:   0.25,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: "Legging đùi"},
					{OptionName: "Size", Value: "S (Dưới 50KG)"},
				},
			},
			{
				SkuCode:  fmt.Sprintf("LEGGING-5DUI-M-%s", uuid.New().String()[:6]),
				Price:    43000,
				Quantity: 50,
				Weight:   0.27,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: "Legging đùi"},
					{OptionName: "Size", Value: "M (Dưới 54KG)"},
				},
			},
		},
	}

	// Option images - array with index [0], [1], etc.
	optionImages := []OptionImageInfo{
		{OptionName: "Màu Sắc", Value: "Xanh than chữ R", ImagePath: "../image_test/mattrc.jpg"},
		{OptionName: "Màu Sắc", Value: "Trắng 3 cúc", ImagePath: "../image_test/anhthe.png"},
		{OptionName: "Màu Sắc", Value: "Đen 3 cúc", ImagePath: "../image_test/QR.jpg"},
	}

	// Main image
	mainImagePath := "../image_test/QR.jpg"

	// Media images
	mediaImagePaths := []string{
		"../image_test/googleplay.png",
		"../image_test/anhthe.png",
		"../image_test/mattrc.jpg",
	}

	err := createProduct(productData, mainImagePath, mediaImagePaths, optionImages)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	t.Log("Product created successfully!")
}

// TestCreateMultipleProducts tests creating multiple products concurrently
func TestCreateMultipleProducts(t *testing.T) {
	numProducts := 3 // Number of products to create

	for i := 0; i < numProducts; i++ {
		t.Run(fmt.Sprintf("Product_%d", i+1), func(t *testing.T) {
			// t.Parallel() // Run tests in parallel

			productData := ProductRequest{
				Name:                      fmt.Sprintf("Áo thun batch test %d - %s", i+1, uuid.New().String()[:8]),
				Key:                       fmt.Sprintf("ao-thun-batch-%d-%s", i+1, uuid.New().String()[:8]),
				Description:               fmt.Sprintf("Mô tả sản phẩm batch test %d với chất liệu cao cấp", i+1),
				ShortDescription:          fmt.Sprintf("Áo thun batch %d chất lượng cao", i+1),
				BrandID:                   "5c5ed281-4b06-495b-ad35-2295c5fba54d",
				CategoryID:                "6ce3146d-9166-4a8d-8dda-038f39a455ec",
				ShopID:                    "uuid-shop",
				ProductIsPermissionReturn: true,
				ProductIsPermissionCheck:  true,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: "Đỏ"},
					{OptionName: "Màu Sắc", Value: "Xanh"},
					{OptionName: "Size", Value: "S (Dưới 50KG)"},
					{OptionName: "Size", Value: "M (Dưới 54KG)"},
				},
				ProductSKU: []ProductSKU{
					{
						SkuCode:  fmt.Sprintf("SKU-RED-S-%d-%s", i+1, uuid.New().String()[:6]),
						Price:    45000,
						Quantity: 100,
						Weight:   0.25,
						OptionValue: []OptionValue{
							{OptionName: "Màu Sắc", Value: "Đỏ"},
							{OptionName: "Size", Value: "S (Dưới 50KG)"},
						},
					},
					{
						SkuCode:  fmt.Sprintf("SKU-RED-M-%d-%s", i+1, uuid.New().String()[:6]),
						Price:    45000,
						Quantity: 100,
						Weight:   0.27,
						OptionValue: []OptionValue{
							{OptionName: "Màu Sắc", Value: "Đỏ"},
							{OptionName: "Size", Value: "M (Dưới 54KG)"},
						},
					},
					{
						SkuCode:  fmt.Sprintf("SKU-BLUE-S-%d-%s", i+1, uuid.New().String()[:6]),
						Price:    45000,
						Quantity: 100,
						Weight:   0.25,
						OptionValue: []OptionValue{
							{OptionName: "Màu Sắc", Value: "Xanh"},
							{OptionName: "Size", Value: "S (Dưới 50KG)"},
						},
					},
					{
						SkuCode:  fmt.Sprintf("SKU-BLUE-M-%d-%s", i+1, uuid.New().String()[:6]),
						Price:    45000,
						Quantity: 100,
						Weight:   0.27,
						OptionValue: []OptionValue{
							{OptionName: "Màu Sắc", Value: "Xanh"},
							{OptionName: "Size", Value: "M (Dưới 54KG)"},
						},
					},
				},
			}

			optionImages := []OptionImageInfo{
				{OptionName: "Màu Sắc", Value: "Đỏ", ImagePath: "../image_test/mattrc.jpg"},
				{OptionName: "Màu Sắc", Value: "Xanh", ImagePath: "../image_test/anhthe.png"},
			}

			mainImagePath := "../image_test/QR.jpg"
			mediaImagePaths := []string{
				"../image_test/googleplay.png",
				"../image_test/anhthe.png",
			}

			err := createProduct(productData, mainImagePath, mediaImagePaths, optionImages)
			if err != nil {
				t.Errorf("Failed to create product %d: %v", i+1, err)
				return
			}

			t.Logf("Product %d created successfully!", i+1)

			// Add a small delay between requests
			time.Sleep(500 * time.Millisecond)
		})
	}
}

// TestCreateProductWithoutOptionImages tests creating a product without option images
func TestCreateProductWithoutOptionImages(t *testing.T) {
	productData := ProductRequest{
		Name:                      fmt.Sprintf("Áo thun đơn giản %s", uuid.New().String()[:8]),
		Key:                       fmt.Sprintf("ao-thun-simple-%s", uuid.New().String()[:8]),
		Description:               "Áo thun đơn giản không có ảnh option, chỉ có 1 size",
		ShortDescription:          "Áo thun simple one size",
		BrandID:                   "5c5ed281-4b06-495b-ad35-2295c5fba54d",
		CategoryID:                "6ce3146d-9166-4a8d-8dda-038f39a455ec",
		ShopID:                    "uuid-shop",
		ProductIsPermissionReturn: true,
		ProductIsPermissionCheck:  true,
		OptionValue: []OptionValue{
			{OptionName: "Size", Value: "Free Size"},
		},
		ProductSKU: []ProductSKU{
			{
				SkuCode:  fmt.Sprintf("SKU-FREESIZE-%s", uuid.New().String()[:6]),
				Price:    35000,
				Quantity: 200,
				Weight:   0.25,
				OptionValue: []OptionValue{
					{OptionName: "Size", Value: "Free Size"},
				},
			},
		},
	}

	mainImagePath := "../image_test/anhthe.png"
	mediaImagePaths := []string{
		"../image_test/QR.jpg",
		"../image_test/googleplay.png",
	}

	// No option images
	err := createProduct(productData, mainImagePath, mediaImagePaths, []OptionImageInfo{})
	if err != nil {
		t.Fatalf("Failed to create simple product: %v", err)
	}

	t.Log("Simple product created successfully!")
}

// createProduct sends a POST request to create a product
func createProduct(productData ProductRequest, mainImagePath string, mediaImagePaths []string, optionImages []OptionImageInfo) error {
	// Create a buffer to write our multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add product data as JSON
	productJSON, err := json.Marshal(productData)
	if err != nil {
		return fmt.Errorf("failed to marshal product data: %w", err)
	}

	err = writer.WriteField("product", string(productJSON))
	if err != nil {
		return fmt.Errorf("failed to write product field: %w", err)
	}

	fmt.Printf("Product JSON: %s\n\n", string(productJSON))

	// Add main image
	if mainImagePath != "" {
		err = addFileToForm(writer, "image", mainImagePath)
		if err != nil {
			return fmt.Errorf("failed to add main image: %w", err)
		}
		fmt.Printf("Added main image: %s\n", mainImagePath)
	}

	// Add media images as array
	for _, mediaPath := range mediaImagePaths {
		err = addFileToForm(writer, "media", mediaPath)
		if err != nil {
			return fmt.Errorf("failed to add media image %s: %w", mediaPath, err)
		}
		fmt.Printf("Added media image: %s\n", mediaPath)
	}

	// Add option images with indexed field names: option_value_images[0], option_value_images[1], etc.
	for i, optImg := range optionImages {
		if optImg.ImagePath != "" {
			// Create field name with index: option_value_images[0], option_value_images[1], etc.
			fieldName := fmt.Sprintf("option_value_images[%d]", i)

			err = addFileToForm(writer, fieldName, optImg.ImagePath)
			if err != nil {
				return fmt.Errorf("failed to add option image %s: %w", optImg.ImagePath, err)
			}
			fmt.Printf("Added option image [%d]: %s for %s = %s\n", i, optImg.ImagePath, optImg.OptionName, optImg.Value)
		}
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", baseURL, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Send request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	fmt.Printf("\n=== Sending request to create product ===\n")
	fmt.Printf("Product: %s\n", productData.Name)
	fmt.Printf("URL: %s\n", baseURL)
	fmt.Printf("Content-Type: %s\n\n", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	fmt.Printf("✅ Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n\n", string(responseBody))

	return nil
}

// addFileToForm adds a file to the multipart form
func addFileToForm(writer *multipart.Writer, fieldName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// TestCreateProductWithManyOptions tests creating a product with many color options
func TestCreateProductWithManyOptions(t *testing.T) {
	productData := ProductRequest{
		Name:                      fmt.Sprintf("Áo thun đa màu %s", uuid.New().String()[:8]),
		Key:                       fmt.Sprintf("ao-thun-multi-color-%s", uuid.New().String()[:8]),
		Description:               "Áo thun với nhiều màu sắc và kích thước, mỗi màu có ảnh riêng",
		ShortDescription:          "Áo thun đa màu sắc",
		BrandID:                   "5c5ed281-4b06-495b-ad35-2295c5fba54d",
		CategoryID:                "6ce3146d-9166-4a8d-8dda-038f39a455ec",
		ShopID:                    "uuid-shop",
		ProductIsPermissionReturn: true,
		ProductIsPermissionCheck:  true,
		OptionValue: []OptionValue{
			{OptionName: "Màu Sắc", Value: "Đỏ"},
			{OptionName: "Màu Sắc", Value: "Xanh"},
			{OptionName: "Màu Sắc", Value: "Vàng"},
			{OptionName: "Màu Sắc", Value: "Trắng"},
			{OptionName: "Size", Value: "S"},
			{OptionName: "Size", Value: "M"},
			{OptionName: "Size", Value: "L"},
		},
		ProductSKU: []ProductSKU{},
	}

	// Generate all combinations of colors and sizes
	colors := []string{"Đỏ", "Xanh", "Vàng", "Trắng"}
	sizes := []string{"S", "M", "L"}

	for _, color := range colors {
		for _, size := range sizes {
			sku := ProductSKU{
				SkuCode:  fmt.Sprintf("SKU-%s-%s-%s", color, size, uuid.New().String()[:6]),
				Price:    50000,
				Quantity: 30,
				Weight:   0.3,
				OptionValue: []OptionValue{
					{OptionName: "Màu Sắc", Value: color},
					{OptionName: "Size", Value: size},
				},
			}
			productData.ProductSKU = append(productData.ProductSKU, sku)
		}
	}

	// Option images for colors
	optionImages := []OptionImageInfo{
		{OptionName: "Màu Sắc", Value: "Đỏ", ImagePath: "../image_test/mattrc.jpg"},
		{OptionName: "Màu Sắc", Value: "Xanh", ImagePath: "../image_test/anhthe.png"},
		{OptionName: "Màu Sắc", Value: "Vàng", ImagePath: "../image_test/googleplay.png"},
		{OptionName: "Màu Sắc", Value: "Trắng", ImagePath: "../image_test/QR.jpg"},
	}

	mainImagePath := "../image_test/QR.jpg"
	mediaImagePaths := []string{
		"../image_test/googleplay.png",
		"../image_test/anhthe.png",
		"../image_test/mattrc.jpg",
	}

	err := createProduct(productData, mainImagePath, mediaImagePaths, optionImages)
	if err != nil {
		t.Fatalf("Failed to create product with many options: %v", err)
	}

	t.Log("Product with many options created successfully!")
	t.Logf("Total SKUs created: %d", len(productData.ProductSKU))
}

// BenchmarkCreateProduct benchmarks the product creation
func BenchmarkCreateProduct(b *testing.B) {
	productData := ProductRequest{
		Name:                      "Benchmark Product",
		Key:                       "benchmark-product",
		Description:               "Benchmark test product",
		ShortDescription:          "Benchmark",
		BrandID:                   "5c5ed281-4b06-495b-ad35-2295c5fba54d",
		CategoryID:                "6ce3146d-9166-4a8d-8dda-038f39a455ec",
		ShopID:                    "uuid-shop",
		ProductIsPermissionReturn: true,
		ProductIsPermissionCheck:  true,
		OptionValue: []OptionValue{
			{OptionName: "Size", Value: "M"},
		},
		ProductSKU: []ProductSKU{
			{
				SkuCode:  "SKU-BENCH-001",
				Price:    40000,
				Quantity: 50,
				Weight:   0.25,
				OptionValue: []OptionValue{
					{OptionName: "Size", Value: "M"},
				},
			},
		},
	}

	mainImagePath := "../image_test/QR.jpg"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		productData.Key = fmt.Sprintf("benchmark-product-%d", i)
		productData.ProductSKU[0].SkuCode = fmt.Sprintf("SKU-BENCH-%d", i)

		err := createProduct(productData, mainImagePath, []string{}, []OptionImageInfo{})
		if err != nil {
			b.Fatalf("Failed to create product: %v", err)
		}
	}
}
