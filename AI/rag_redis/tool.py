"""
Tool to fetch product data and index into vector search API.
Optimized with Multi-threading using ThreadPoolExecutor.
"""

import requests
import time
import threading
from typing import List
from concurrent.futures import ThreadPoolExecutor, as_completed

# --- CẤU HÌNH ---
PRODUCT_API_BASE = "http://172.26.127.95:9001/v1/product"
VECTOR_API_BASE = "http://localhost:9101"

GET_ALL_PRODUCTS_URL = f"{PRODUCT_API_BASE}/getallproductid"
BUILD_SEARCH_STRING_URL = f"{PRODUCT_API_BASE}/build_search_string"
ADD_DOCUMENT_URL = f"{VECTOR_API_BASE}/documents"

MAX_WORKERS = 15  # Số luồng chạy song song (bạn có thể tăng/giảm tùy ý)
print_lock = threading.Lock()  # Khóa để in log không bị lỗi hiển thị khi chạy song song

# --- CÁC HÀM API ---

def fetch_all_product_ids() -> List[str]:
    """Lấy danh sách tất cả Product ID."""
    print("[INFO] Fetching all product IDs...")
    try:
        response = requests.get(GET_ALL_PRODUCTS_URL, timeout=30)
        response.raise_for_status()
        data = response.json()
        
        if data.get("code") == 200 and data.get("status") == "success":
            product_ids = data.get("result", {}).get("product_ids", [])
            print(f"[SUCCESS] Retrieved {len(product_ids)} product IDs")
            return product_ids
        else:
            print(f"[ERROR] Failed to fetch product IDs: {data.get('message')}")
            return []
    except Exception as e:
        print(f"[ERROR] Exception while fetching product IDs: {e}")
        return []


def build_search_string(product_id: str) -> str:
    """Gọi API để tạo chuỗi search string cho 1 sản phẩm."""
    url = f"{BUILD_SEARCH_STRING_URL}/{product_id}"
    try:
        response = requests.get(url, timeout=30)
        response.raise_for_status()
        data = response.json()
        
        if data.get("code") == 200 and data.get("status") == "success":
            return data.get("result", {}).get("search_string", "")
        else:
            with print_lock:
                print(f"[WARNING] Failed to build search string for {product_id}: {data.get('message')}")
            return ""
    except Exception as e:
        with print_lock:
            print(f"[ERROR] Exception while building search string for {product_id}: {e}")
        return ""


def index_document(doc_id: str, text_content: str) -> bool:
    """Gọi API để index dữ liệu vào Vector DB."""
    payload = {
        "doc_id": doc_id,
        "text_content": text_content,
        "doc_type":"product"
    }
    
    try:
        response = requests.post(ADD_DOCUMENT_URL, json=payload, timeout=60)
        response.raise_for_status()
        data = response.json()
        
        if data.get("status") == "success":
            return True
        else:
            with print_lock:
                print(f"[WARNING] Failed to index document {doc_id}: {data}")
            return False
    except Exception as e:
        with print_lock:
            print(f"[ERROR] Exception while indexing document {doc_id}: {e}")
        return False


def process_single_product(product_id: str) -> bool:
    """
    Hàm xử lý cho 1 sản phẩm (sẽ được chạy bởi 1 luồng riêng biệt).
    Luồng xử lý: Lấy chuỗi search -> Index
    """
    # 1. Build search string
    search_string = build_search_string(product_id)
    
    if not search_string:
        with print_lock:
            print(f"  [SKIP] {product_id}: Empty search string")
        return False
    
    # 2. Index into vector API
    if index_document(product_id, search_string):
        with print_lock:
            print(f"  [SUCCESS] {product_id} (Len: {len(search_string)})")
        return True
    else:
        with print_lock:
            print(f"  [FAILED] {product_id}: Could not index")
        return False


# --- HÀM MAIN ---

def main():
    print("=" * 80)
    print(f"Starting Product Indexing Tool (Multi-threaded: {MAX_WORKERS} workers)")
    print("=" * 80)
    
    # B1: Lấy danh sách ID
    product_ids = fetch_all_product_ids()
    
    if not product_ids:
        print("[ERROR] No product IDs found. Exiting.")
        return
    
    total = len(product_ids)
    success_count = 0
    failed_count = 0
    
    print(f"\n[INFO] Starting to process {total} products...\n")
    
    # B2: Sử dụng ThreadPoolExecutor để chạy song song
    with ThreadPoolExecutor(max_workers=MAX_WORKERS) as executor:
        # Submit các task vào pool
        future_to_pid = {executor.submit(process_single_product, pid): pid for pid in product_ids}
        
        # Nhận kết quả khi từng task hoàn thành
        for i, future in enumerate(as_completed(future_to_pid), 1):
            try:
                result = future.result()
                if result:
                    success_count += 1
                else:
                    failed_count += 1
                
                # In thông báo tiến độ mỗi khi xong 10 sản phẩm
                if i % 10 == 0 or i == total:
                    print(f"--- Progress: {i}/{total} completed ---")
                    
            except Exception as exc:
                print(f"[CRITICAL ERROR] Thread exception: {exc}")
                failed_count += 1

    # Tổng kết
    print("\n" + "=" * 80)
    print("Indexing Summary")
    print("=" * 80)
    print(f"Total products: {total}")
    print(f"Successfully indexed: {success_count}")
    print(f"Failed: {failed_count}")
    print("=" * 80)


if __name__ == "__main__":
    main()