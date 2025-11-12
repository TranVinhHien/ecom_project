import os
import requests
import json
from typing import List, Dict, Any, Optional
BASE_URL = "https://lemarchenoble.id.vn/api/v1/profile/users/profiles/get-my-profile"


def get_user_info(token: str) -> Optional[Dict[str, Any]]:
    """
    Calls the API endpoint to get user information.

    Args:
        token (str): Authentication token to include in the request header

    Returns:
        Optional[Dict[str, Any]]: A dictionary containing user information, or None if not found

    Raises:
        requests.RequestException: If there's an error with the API request
        ValueError: If the API returns invalid data
    """
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    try:
        response = requests.get(BASE_URL, headers=headers)
        response.raise_for_status()

        data = response.json()
        if data.get("code") != 10000:
            raise ValueError("API response indicates failure or invalid data")
        return data.get("result")

    except requests.RequestException as e:
        print(f"Error calling user info API: {str(e)}")
        raise
    except (json.JSONDecodeError, ValueError) as e:
        print(f"Error parsing API response: {str(e)}")
        raise
