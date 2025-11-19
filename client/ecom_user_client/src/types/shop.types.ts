export interface TaxInfo {
  id: string;
  taxCode: string;
  taxNationalName: string;
  taxShortName: string;
  taxPresentName: string;
  taxActiveDate: string;
  taxBusinessType: string;
  taxActiveStatus: boolean;
}

export interface ShopInfo {
  taxInfo: TaxInfo;
  id: string;
  shopName: string;
  shopDescription: string;
  shopLogo: string;
  shopAddress: string;
  shopPersonalIdentifyId: string;
  shopEmail: string;
  shopPhone: string;
  shopStatus: boolean;
  walletAmount: number;
  followerCount: number;
  isFollowing: boolean;
  createdDate: string;
}

export interface ShopApiResponse {
  result: ShopInfo;
  messages: string[];
  succeeded: boolean;
  code: number;
}
