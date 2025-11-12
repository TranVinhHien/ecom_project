import os
import requests
import json
from typing import List, Dict, Any, Optional
BASE_URL = "http://172.26.127.95:9001/v1"
ROUTER= {
    "product_search": f"{BASE_URL}/product/get_products_detail_for_search",
}

# def get_product_detail_for_search(product_ids: Dict[str, float]) -> Optional[Dict[str, Any]]:
#     """
#     Calls the API endpoint to get product details for a list of product IDs.

#     Args:
#         product_ids (Dict[str, float]): A dictionary of product IDs to retrieve details for

#     Returns:
#         Optional[Dict[str, Any]]: A dictionary containing product details, or None if not found
#     """
#     headers = {
#         "Content-Type": "application/json"
#     }
#     doc_ids = list(product_ids.keys())
#     try:
#         response = requests.get(ROUTER.get("product_search"), headers=headers, params={"product_ids": doc_ids})
#         response.raise_for_status()
#         data = response.json()
#         if data.get("code") != 200:
#             raise ValueError("API response indicates failure or invalid data")
#         product = data.get("result").get("data")
        
#         return 

#     except requests.RequestException as e:
#         print(f"Error calling product detail API: {str(e)}")
#         raise
#     except (json.JSONDecodeError, ValueError) as e:
#         print(f"Error parsing API response: {str(e)}")
#         raise

def get_product_detail_for_search(product_ids: Dict[str, float]) -> Optional[Dict[str, Any]]:
    """
    Calls the API endpoint to get product details and maps similarity scores.

    Args:
        product_ids (Dict[str, float]): Một dictionary với key là product_id
                                          và value là similarity_score.

    Returns:
        Optional[Dict[str, Any]]: Một dictionary với key là product_id và 
                                   value là object product (đã bao gồm 
                                   similarity_score), hoặc None nếu lỗi.
    """
    headers = {
        "Content-Type": "application/json"
    }
    
    # Lấy danh sách các ID để gửi cho API
    doc_ids = list(product_ids.keys())
    
    # Nếu không có ID nào cần tìm, trả về dict rỗng
    if not doc_ids:
        print("[INFO] No product IDs to search.")
        return {}

    try:
        # Gửi request chỉ với danh sách ID
        response = requests.get(
            ROUTER.get("product_search"), 
            headers=headers, 
            params={"product_ids": doc_ids}
        )
        response.raise_for_status()  # Báo lỗi nếu status code là 4xx hoặc 5xx
        data = response.json()

        if data.get("code") != 200:
            print(f"[ERROR] API response indicates failure: {data.get('message')}")
            raise ValueError("API response indicates failure or invalid data")

        # Lấy danh sách sản phẩm từ API
        products_list: List[Dict[str, Any]] = data.get("result", {}).get("data")

        if not isinstance(products_list, list):
            print("[ERROR] API response 'data' is not a list.")
            return {}

        # --- THAY ĐỔI QUAN TRỌNG: Map điểm similarity ---
        
        # Tạo một dictionary mới để lưu kết quả
        final_products_map: List[str] = []
        
        for item in products_list:
            # Lấy ID sản phẩm từ trong cấu trúc JSON trả về
            product_data = item.get("product", {})
            doc_id = product_data.get("id")

            # Nếu ID này có trong dict 'product_ids' đầu vào của chúng ta
            if doc_id and doc_id in product_ids:
                
                # 1. Lấy điểm similarity_score từ dict đầu vào
                similarity = product_ids[doc_id]
                
                # 2. Thêm trường 'similarity_score' vào object 'item'
                item["similarity_score"] = similarity
                
                # 3. Thêm 'item' (đã có score) vào dict kết quả
                #    dùng doc_id làm key
                final_products_map.append(item)
        final_products_map.sort(key=lambda item: item.get('similarity_score', 0.0), reverse=True)
        # Trả về dict đã map
        return final_products_map

    except requests.exceptions.RequestException as e:
        print(f"[ERROR] HTTP Request failed: {e}")
        return None  # Hoặc {} tùy vào cách bạn xử lý lỗi
    except ValueError as e:
        print(f"[ERROR] API Data error: {e}")
        return None  # Hoặc {}
    except Exception as e:
        print(f"[ERROR] An unexpected error occurred: {e}")
        return None  # Hoặc {}