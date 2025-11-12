"""
Tool to fetch product data and index into vector search API.

This script:
1. Fetches all product IDs from the product service
2. For each product ID, builds a search string
3. Indexes each product into the local vector search API
"""

import requests
import time
from typing import List, Dict, Any


# API endpoints
PRODUCT_API_BASE = "http://172.26.127.95:9001/v1/product"
VECTOR_API_BASE = "http://localhost:9101"

GET_ALL_PRODUCTS_URL = f"{PRODUCT_API_BASE}/getallproductid"
BUILD_SEARCH_STRING_URL = f"{PRODUCT_API_BASE}/build_search_string"
ADD_DOCUMENT_URL = f"{VECTOR_API_BASE}/documents"


def fetch_all_product_ids() -> List[str]:
    """Fetch all product IDs from the product service."""
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
    """Build search string for a specific product ID."""
    url = f"{BUILD_SEARCH_STRING_URL}/{product_id}"
    try:
        response = requests.get(url, timeout=30)
        response.raise_for_status()
        data = response.json()
        
        if data.get("code") == 200 and data.get("status") == "success":
            search_string = data.get("result", {}).get("search_string", "")
            return search_string
        else:
            print(f"[WARNING] Failed to build search string for {product_id}: {data.get('message')}")
            return ""
    except Exception as e:
        print(f"[ERROR] Exception while building search string for {product_id}: {e}")
        return ""


def index_document(doc_id: str, text_content: str) -> bool:
    """Index a document into the vector search API."""
    payload = {
        "doc_id": doc_id,
        "text_content": text_content
    }
    
    try:
        response = requests.post(ADD_DOCUMENT_URL, json=payload, timeout=60)
        response.raise_for_status()
        data = response.json()
        
        if data.get("status") == "success":
            return True
        else:
            print(f"[WARNING] Failed to index document {doc_id}: {data}")
            return False
    except Exception as e:
        print(f"[ERROR] Exception while indexing document {doc_id}: {e}")
        return False


def main():
    """Main execution flow."""
    print("=" * 80)
    print("Starting Product Indexing Tool")
    print("=" * 80)
    
    # Step 1: Fetch all product IDs
    product_ids =fetch_all_product_ids()
    
    if not product_ids:
        print("[ERROR] No product IDs found. Exiting.")
        return
    
    # Step 2 & 3: For each product, build search string and index
    total = len(product_ids)
    success_count = 0
    failed_count = 0
    
    print(f"\n[INFO] Starting to process {total} products...\n")
    
    for idx, product_id in enumerate(product_ids, 1):
        print(f"[{idx}/{total}] Processing product: {product_id}")
        
        # Build search string
        search_string = build_search_string(product_id)
        
        if not search_string:
            print(f"  └─ [SKIP] Empty search string for {product_id}")
            failed_count += 1
            continue
        
        print(f"  └─ [OK] Search string length: {len(search_string)} characters")
        
        # Index into vector API
        if index_document(product_id, search_string):
            print(f"  └─ [SUCCESS] Indexed product {product_id}")
            success_count += 1
        else:
            print(f"  └─ [FAILED] Could not index product {product_id}")
            failed_count += 1
        
        # Small delay to avoid overwhelming the APIs
        time.sleep(0.1)
        print()
    
    # Summary
    print("=" * 80)
    print("Indexing Summary")
    print("=" * 80)
    print(f"Total products: {total}")
    print(f"Successfully indexed: {success_count}")
    print(f"Failed: {failed_count}")
    print("=" * 80)


if __name__ == "__main__":
    main()
