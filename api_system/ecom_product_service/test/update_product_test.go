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
)

const (
	updateProductURL = "http://172.26.127.95:9001/v1/product/update/%s"
	updateToken      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJoaWVubGF6YWRhIiwic2NvcGUiOiJTWVNURU0iLCJpc3MiOiJsZW1hcmNoZW5vYmxlLmlkLnZuIiwiZXhwIjo0OTE3NTExMjUyLCJpYXQiOjE3NjE3NTEyNTIsInVzZXJJZCI6IjE2NzQwOGUzLWFmZWYtNDhiOS04ZTRmLTZkZDQxZWJmMzQ2NCIsImp0aSI6ImU2YzgyN2E2LTIyOTYtNGNlOC1iMjQ1LWM3MDIxNWM4MGJjNyIsImVtYWlsIjoidmluaGhpZW4xMnpAZ21haWwuY29tIn0.CPnP_NqB_WtaQb9X43YKFav8wYzdqB14jFNtnPr74as"
	testProductID    = "09047bfd-dd5d-415c-a656-3638991d03a4"
	testOptionID1    = "36b1155d-8b26-41e2-9325-2c04cc430907"
	testOptionID2    = "47c094f8-d2d7-4d14-a333-690163a586cf"
	testOptionID3    = "6587c8ef-05f5-460e-abe9-8dd6b201058f"
)

type ProductUpdateTest struct {
	Name             *string      `json:"name,omitempty"`
	Description      *string      `json:"description,omitempty"`
	BrandID          *string      `json:"brand_id,omitempty"`
	CategoryID       *string      `json:"category_id,omitempty"`
	UrlImage         *string      `json:"url_image,omitempty"`
	UrlMedia         *string      `json:"url_media,omitempty"`
	Options          []OptionTest `json:"options,omitempty"`
	IsHaveOptionSKUs *bool        `json:"is_have_option_skus,omitempty"`
}

type OptionTest struct {
	ID          *string           `json:"id,omitempty"`
	Name        *string           `json:"name,omitempty"`
	OptionValue []OptionValueTest `json:"option_value,omitempty"`
}

type OptionValueTest struct {
	ID     *string `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	UrlImg *string `json:"url_img,omitempty"`
}

type MediaItemTest struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	TypeFile string `json:"type_file"`
}

func stringPtrUpdate(s string) *string {
	return &s
}

func boolPtrUpdate(b bool) *bool {
	return &b
}

func sendUpdateRequest(productID string, product ProductUpdateTest, imageFile string, mediaFiles []string, optionImages map[int]string) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	productJSON, err := json.Marshal(product)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal product: %v", err)
	}
	if err := writer.WriteField("product", string(productJSON)); err != nil {
		return nil, fmt.Errorf("failed to write product field: %v", err)
	}

	if imageFile != "" {
		file, err := os.Open(imageFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open image file: %v", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("image", filepath.Base(imageFile))
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %v", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			return nil, fmt.Errorf("failed to copy file: %v", err)
		}
	}

	if len(mediaFiles) > 0 {
		mediaArray := make([]MediaItemTest, 0)
		for _, mediaFile := range mediaFiles {
			file, err := os.Open(mediaFile)
			if err != nil {
				return nil, fmt.Errorf("failed to open media file: %v", err)
			}
			defer file.Close()

			part, err := writer.CreateFormFile("media", filepath.Base(mediaFile))
			if err != nil {
				return nil, fmt.Errorf("failed to create media form file: %v", err)
			}
			if _, err := io.Copy(part, file); err != nil {
				return nil, fmt.Errorf("failed to copy media file: %v", err)
			}

			mediaArray = append(mediaArray, MediaItemTest{
				Name:     filepath.Base(mediaFile),
				URL:      "",
				TypeFile: "image",
			})
		}

		mediaJSON, err := json.Marshal(mediaArray)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal media: %v", err)
		}
		if err := writer.WriteField("media", string(mediaJSON)); err != nil {
			return nil, fmt.Errorf("failed to write media field: %v", err)
		}
	}

	for idx, optionImageFile := range optionImages {
		file, err := os.Open(optionImageFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open option image file: %v", err)
		}
		defer file.Close()

		fieldName := fmt.Sprintf("option_value_images[%d]", idx)
		part, err := writer.CreateFormFile(fieldName, filepath.Base(optionImageFile))
		if err != nil {
			return nil, fmt.Errorf("failed to create option image form file: %v", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			return nil, fmt.Errorf("failed to copy option image file: %v", err)
		}
	}

	writer.Close()

	url := fmt.Sprintf(updateProductURL, productID)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+updateToken)

	client := &http.Client{}
	return client.Do(req)
}

func TestUpdateProductName(t *testing.T) {
	fmt.Println("\n========== TEST 1: CẬP NHẬT CHỈ TÊN SẢN PHẨM ==========")

	newName := "Áo Thun Nam Cập Nhật - Test 1"
	product := ProductUpdateTest{
		Name: stringPtrUpdate(newName),
	}

	resp, err := sendUpdateRequest(testProductID, product, "", nil, nil)
	if err != nil {
		t.Fatalf("❌ Lỗi khi gửi request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("❌ Mã trạng thái không mong đợi: %d", resp.StatusCode)
	} else {
		fmt.Println("✅ CẬP NHẬT TÊN SẢN PHẨM THÀNH CÔNG")
	}
}

func TestUpdateProductWithImage(t *testing.T) {
	fmt.Println("\n========== TEST 2: CẬP NHẬT SẢN PHẨM VỚI ẢNH MỚI ==========")

	newName := "Áo Thun Nam - Cập Nhật Với Ảnh Mới"
	newDesc := "Mô tả sản phẩm đã được cập nhật với ảnh mới"
	product := ProductUpdateTest{
		Name:        stringPtrUpdate(newName),
		Description: stringPtrUpdate(newDesc),
	}

	imageFile := "../image_test/anh1.txt"
	resp, err := sendUpdateRequest(testProductID, product, imageFile, nil, nil)
	if err != nil {
		t.Fatalf("❌ Lỗi khi gửi request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("❌ Mã trạng thái không mong đợi: %d", resp.StatusCode)
	} else {
		fmt.Println("✅ CẬP NHẬT SẢN PHẨM VỚI ẢNH MỚI THÀNH CÔNG")
	}
}

func TestUpdateProductWithMedia(t *testing.T) {
	fmt.Println("\n========== TEST 3: CẬP NHẬT SẢN PHẨM VỚI MEDIA FILES ==========")

	newName := "Áo Thun Nam - Cập Nhật Với Media"
	product := ProductUpdateTest{
		Name: stringPtrUpdate(newName),
	}

	imageFile := "../image_test/anh1.txt"
	mediaFiles := []string{
		"../image_test/anh1.txt",
	}

	resp, err := sendUpdateRequest(testProductID, product, imageFile, mediaFiles, nil)
	if err != nil {
		t.Fatalf("❌ Lỗi khi gửi request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("❌ Mã trạng thái không mong đợi: %d", resp.StatusCode)
	} else {
		fmt.Println("✅ CẬP NHẬT SẢN PHẨM VỚI MEDIA THÀNH CÔNG")
	}
}

func TestUpdateProductOptions(t *testing.T) {
	fmt.Println("\n========== TEST 4: CẬP NHẬT OPTIONS SẢN PHẨM ==========")

	product := ProductUpdateTest{
		Options: []OptionTest{
			{
				ID:   stringPtrUpdate(testOptionID1),
				Name: stringPtrUpdate("Màu sắc"),
				OptionValue: []OptionValueTest{
					{
						ID:   stringPtrUpdate(testOptionID2),
						Name: stringPtrUpdate("Đỏ"),
					},
					{
						ID:   stringPtrUpdate(testOptionID3),
						Name: stringPtrUpdate("Xanh"),
					},
					{
						Name: stringPtrUpdate("Vàng"),
					},
				},
			},
		},
	}

	optionImages := map[int]string{
		0: "../image_test/anh1.txt",
		1: "../image_test/anh1.txt",
		2: "../image_test/anh1.txt",
	}

	resp, err := sendUpdateRequest(testProductID, product, "", nil, optionImages)
	if err != nil {
		t.Fatalf("❌ Lỗi khi gửi request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("❌ Mã trạng thái không mong đợi: %d", resp.StatusCode)
	} else {
		fmt.Println("✅ CẬP NHẬT OPTIONS SẢN PHẨM THÀNH CÔNG")
	}
}

func TestUpdateProductComplete(t *testing.T) {
	fmt.Println("\n========== TEST 5: CẬP NHẬT TOÀN BỘ THÔNG TIN SẢN PHẨM ==========")

	newName := "Áo Thun Nam Cao Cấp - Cập Nhật Hoàn Chỉnh"
	newDesc := "Áo thun nam chất liệu cotton 100%, thoáng mát, thấm hút mồ hôi tốt. Đã được cập nhật hoàn toàn."
	brandID := "01932d39-3d7c-728c-850c-1a60ed7ed40a"
	categoryID := "01932d39-41a9-757f-a5e8-27e36c19bd2b"

	product := ProductUpdateTest{
		Name:        stringPtrUpdate(newName),
		Description: stringPtrUpdate(newDesc),
		BrandID:     stringPtrUpdate(brandID),
		CategoryID:  stringPtrUpdate(categoryID),
		Options: []OptionTest{
			{
				ID:   stringPtrUpdate(testOptionID1),
				Name: stringPtrUpdate("Màu sắc"),
				OptionValue: []OptionValueTest{
					{
						ID:   stringPtrUpdate(testOptionID2),
						Name: stringPtrUpdate("Đỏ Đậm"),
					},
					{
						ID:   stringPtrUpdate(testOptionID3),
						Name: stringPtrUpdate("Xanh Navy"),
					},
				},
			},
		},
		IsHaveOptionSKUs: boolPtrUpdate(true),
	}

	imageFile := "../image_test/anh1.txt"
	mediaFiles := []string{
		"../image_test/anh1.txt",
	}
	optionImages := map[int]string{
		0: "../image_test/anh1.txt",
		1: "../image_test/anh1.txt",
	}

	resp, err := sendUpdateRequest(testProductID, product, imageFile, mediaFiles, optionImages)
	if err != nil {
		t.Fatalf("❌ Lỗi khi gửi request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("❌ Mã trạng thái không mong đợi: %d", resp.StatusCode)
	} else {
		fmt.Println("✅ CẬP NHẬT TOÀN BỘ THÔNG TIN SẢN PHẨM THÀNH CÔNG")
	}
}
