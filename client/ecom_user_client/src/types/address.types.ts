export interface Province {
  id: string;
  name: string;
  fullName: string;
  type: string;
}

export interface District {
  id: string;
  name: string;
  fullName: string;
  type: string;
  province: Province;
}

export interface Ward {
  id: string;
  name: string;
  fullName: string;
  type: string;
  district: District;
}

export interface AddressFormData {
  name: string;
  phoneNumber: string;
  address: {
    wardId: string;
    other: string;
  };
}

export interface UserAddress {
  id: string;
  name: string;
  phoneNumber: string;
  address: {
    id: string;
    ward: Ward;
    other: string;
  };
}
