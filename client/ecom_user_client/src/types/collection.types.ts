

export type event_type =  'search' | 'cart_add' | 'cart_remove' | 'wishlist' | 'purchase'| 'click';

export interface CollectionType {
  user_id: string;
  event_type: event_type;
  product_id?: string;
  shop_id?: string;
  price?: number;
  quantity?: number;
  metadata?: Record<string, any>;
}

export interface CollectionApiResponse {
  status: string;

}
